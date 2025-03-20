package email

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kidusshun/ecom_bot/utils"
)

type Handler struct {
	service EmailService
}

func NewHandler(service EmailService) *Handler {
	return &Handler{
		service: service,
	}
}


func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/track-activity", h.sendEmail)
}

func (h *Handler) sendEmail(w http.ResponseWriter, r *http.Request) {

	var request SendEmailPayload

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Println(request)

	err = h.service.SendEmailService(request.Email)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Email sent successfully"})

}