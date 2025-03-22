package email

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/config"
	"github.com/kidusshun/ecom_bot/llmclient"
	"github.com/kidusshun/ecom_bot/service/chat"
	"github.com/kidusshun/ecom_bot/service/user"
)

type Service struct {
	client llmclient.LlmClient
	store chat.ChatStore
	userStore user.UserStore
}


func NewService(store chat.ChatStore, client llmclient.LlmClient, userStore user.UserStore) *Service {
	return &Service{
		client: client,
		store:store,
		userStore: userStore,
	}
}

type EmailMessage struct {
	To          string
	Subject     string
	Body        string
	ContentType string // "text/plain" or "text/html"
}

func (s *Service) SendEmailService(email string) error {
	messageHistory, err := s.store.GetChatHistory(uuid.New())
	
	if err != nil {
		return err
	}
	response, err := s.client.GenerateEmail(messageHistory)

	
	if err != nil {
		return err
	}
	
	results := response.Candidates[0].Content.Parts[0].Text
	log.Println("response from gemini",results)
	var responseMap map[string]string
	err = json.Unmarshal([]byte(results), &responseMap)

	if err != nil {
		return err
	}

	log.Println(responseMap)

	return sendEmail(responseMap, email)

}


func (s *Service) SendCouponEmail(coupon CouponRequest) error {
	response, err := s.client.GenerateCouponEmail(coupon.Code, coupon.Discount, coupon.ExpirationDate)

	if err != nil {
		return err
	}
	
	results := response.Candidates[0].Content.Parts[0].Text
	log.Println("response from gemini",results)
	var responseMap map[string]string
	err = json.Unmarshal([]byte(results), &responseMap)

	if err != nil {
		return err
	}

	log.Println(responseMap)

	users, err := s.userStore.GetAllUsers()

	
	if err != nil {
		return err
	}
	log.Println("got users")

	err = nil
	for _, user := range *users {
		err =sendEmail(responseMap, user.Email)

		log.Println("sent email to", user.Email, err)
	}

	return err
}


func sendEmail(responseMap map[string]string, email string) error {
	auth := smtp.PlainAuth("", config.SMTPEnvs.SenderEmail, config.SMTPEnvs.AppPassword, config.SMTPEnvs.SMTPServer)
	message := EmailMessage{
		To: email,
		Subject: responseMap["subject"],
		Body: responseMap["body"],
		ContentType: "text/plain",
	}

	headers := make(map[string]string)
	headers["From"] = config.SMTPEnvs.SenderEmail
	headers["To"] = message.To
	headers["Subject"] = message.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = message.ContentType + "; charset=\"utf-8\""

	var emailBody strings.Builder
	for key, value := range headers {
		emailBody.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	emailBody.WriteString("\r\n")
	emailBody.WriteString(message.Body)
	
	// Connect to the SMTP server and send the email
	addr := fmt.Sprintf("%s:%s", config.SMTPEnvs.SMTPServer, config.SMTPEnvs.SMTPPort)
	return smtp.SendMail(
		addr,
		auth,
		config.SMTPEnvs.SenderEmail,
		[]string{message.To},
		[]byte(emailBody.String()),
	)
}