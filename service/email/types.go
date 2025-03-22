package email

import "github.com/google/uuid"

type EmailService interface {
	SendEmailService(email string) error
	SendCouponEmail(coupon CouponRequest) error
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

type CouponRequest struct {
	ID            uuid.UUID     `json:"id"`
	Code          string  `json:"code"`
	Discount      float64 `json:"discount"`
	ExpirationDate    string  `json:"expiration_date"`
	IsActive bool  `json:"is_active"`
}