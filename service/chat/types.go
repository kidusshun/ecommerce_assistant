package chat

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/llmclient"
	"github.com/kidusshun/ecom_bot/service/product"
)

type ChatStore interface {
	WriteMessage(message string, sessionID uuid.UUID) error
	GetChatHistory(sessionID uuid.UUID) ([]llmclient.Message, error)
}

type ChatService interface {
	Chat(request ChatRequest) (ChatResponse, error)
}

type ChatRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	Message string `json:"message" validate:"required"`
	Attachment []ImageInfo `json:"attachment"`
}
type ImageInfo struct {
	File     *multipart.FileHeader
	MimeType string
	Content  []byte
}

type ChatResponse struct {
	Content string `json:"content"`
	Role string `json:"role"`
	Attachment []ImageInfo `json:"attachment"`
	Location string `json:"location"`
	Products []product.Product `json:"products"`
}