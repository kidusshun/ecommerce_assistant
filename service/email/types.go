package email

type EmailService interface {
	SendEmailService(email string) error
}


type SendEmailPayload struct {
	Email string 	`json:"email"`
}

type SMTPConfig struct {
	SMTPServer   string
	SMTPPort     string
	SenderEmail  string
	AppPassword  string
}