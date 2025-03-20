package chat

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/service/auth"
	"github.com/kidusshun/ecom_bot/utils"
)

type Handler struct {
	service ChatService
}

func NewHandler(service ChatService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.With(auth.CheckBearerToken).Post("/chat", h.handleChat)
	router.With(auth.CheckBearerToken).Get("/chat/history", h.handleChatHistory)
}

func (h *Handler) handleChat(w http.ResponseWriter, r *http.Request) {
	// Parse the request

	
	chatRequest, err := parseRequest(w, r)

	if err != nil {
		utils.WriteError(w, 400, err)
		return
	}

	// Call the service
	response, err := h.service.Chat(chatRequest)
	log.Println("here",response.Content)
	if err != nil {
		log.Println(err)
		utils.WriteError(w, 500, err)
		return
	}
	
	// Return the response
	utils.WriteJSON(w, 200, response)
}

func (h *Handler) handleChatHistory(w http.ResponseWriter, r *http.Request){
	log.Println(MessageHistory)
	utils.WriteJSON(w, http.StatusOK, MessageHistory)
}

func parseRequest(w http.ResponseWriter,r *http.Request) (ChatRequest, error){
	
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err := r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
	if err != nil {
		return ChatRequest{}, err
	}

	var chatRequest ChatRequest
	
	// Extract the sessionID and message from the form data
	sessionIDStr := r.FormValue("session_id")
	message := r.FormValue("message")
	
	if message == "" {
		return ChatRequest{}, err
	}
	chatRequest.Message = message

	if sessionIDStr != "" {
	// Parse the sessionID from string to uuid
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return ChatRequest{}, err
	}
	chatRequest.SessionID = sessionID
	}

	files := r.MultipartForm.File["attachment"]
    for _, fileHeader := range files {
        file, err := fileHeader.Open()
    if err != nil {
        http.Error(w, "Unable to open file", http.StatusInternalServerError)
        return ChatRequest{},err
    }
    defer file.Close()

    // Read the file content
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Unable to read file", http.StatusInternalServerError)
        return ChatRequest{},err
    }

	base64Encoded := base64.StdEncoding.EncodeToString(fileBytes)

    // Populate ImageInfo with metadata and content
    attachment := ImageInfo{
        File:     fileHeader,                 // Metadata about the file
        MimeType: fileHeader.Header.Get("Content-Type"), // MIME type
        Content:  []byte(base64Encoded),                 // File content as a byte array
    }
        chatRequest.Attachment = append(chatRequest.Attachment, attachment)
    }
	return chatRequest, nil
}

func isValidImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}
