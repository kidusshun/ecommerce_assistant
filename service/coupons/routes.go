package coupons

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kidusshun/ecom_bot/utils"
)

type Handler struct {
	service CouponService
}

func NewHandler(service CouponService) *Handler {
	return &Handler{
		service: service,
	}
}


func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Post("/create-coupon", h.createCoupon)
}

func (h *Handler) createCoupon(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var request CreateCouponRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		utils.WriteError(w, 400, err)
		return
	}

	// Call the service
	response, err := h.service.CreateCoupon(request)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// Return the response
	utils.WriteJSON(w, http.StatusOK, response)
}