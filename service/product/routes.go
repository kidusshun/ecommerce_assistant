package product

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kidusshun/ecom_bot/service/auth"
	"github.com/kidusshun/ecom_bot/service/user"
	"github.com/kidusshun/ecom_bot/utils"
)

type Handler struct {
	store     ProductStore
	userStore user.UserStore
}

func NewHandler(store ProductStore, userStore user.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.With(auth.CheckBearerToken).Get("/products", h.handleGetProducts)
    router.With(auth.CheckBearerToken).Get("/products/{productID}", h.handleGetProductByID)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)
	log.Println("emaaaaaaaaaaaaaaail",userEmail)
	user, err := h.userStore.GetUserByEmail(userEmail)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	log.Println("user",user)
	products, err := h.store.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, 200, products)
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	id, err := uuid.Parse(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := h.store.GetProductByID(id)
	if err != nil {
		utils.WriteError(w, 500, err)
		return
	}

	utils.WriteJSON(w, 200, product)

}
