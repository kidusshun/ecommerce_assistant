package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kidusshun/ecom_bot/llmclient"
	"github.com/kidusshun/ecom_bot/service/chat"
	"github.com/kidusshun/ecom_bot/service/email"
	"github.com/kidusshun/ecom_bot/service/product"
	"github.com/kidusshun/ecom_bot/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // Allow specific origin
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        ExposedHeaders:   []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           300, // Max cache duration in seconds
    }))
	router.Use(middleware.Logger)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(router)

	chatStore := chat.NewStore(s.db)
	llmTools := llmclient.NewQueryStore(s.db)
	client := llmclient.NewLlmClient(llmTools)
	chatService := chat.NewChatService(chatStore, client)
	chatHandler := chat.NewHandler(chatService)
	chatHandler.RegisterRoutes(router)

	emailService := email.NewService(chatStore, client)
	emailHandler := email.NewHandler(emailService)
	emailHandler.RegisterRoutes(router)


	log.Println("Listening on ", s.addr)
	err := http.ListenAndServe(s.addr, router)

	return err
}
