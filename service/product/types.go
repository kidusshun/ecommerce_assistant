package product

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	Image         string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductNOImage struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateProductPayload struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	Image         string  `json:"image_url"`
	Price         float64 `json:"price" validate:"required"`
	StockQuantity int     `json:"stock_quantity" validate:"required"`
}

type ProductStore interface {
	GetProducts() (*[]Product, error)
	GetProductByID(id uuid.UUID) (*Product, error)
	ScanRowsIntoProduct(rows *sql.Rows) (*[]Product, error)
	AddProduct(product Product) error
}

