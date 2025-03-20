package product

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetProducts()(*[]Product, error) {
	rows, err := s.db.Query("SELECT id, name, description, price, stock_quanity, image, created_at, updated_at from products")
	if err != nil {
		return nil, err
	}
	products, err := s.ScanRowsIntoProduct(rows)
	if err != nil {
		return nil, err
	}
	return products, nil

}

func (s *Store) AddProduct(product Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, price, stock_quanity, image) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at", product.Name, product.Description, product.Price, product.StockQuantity, product.Image)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProductByID(id uuid.UUID)(*Product, error) {
	row := s.db.QueryRow("SELECT id, name, description, price, stock_quanity, image, created_at, updated_at FROM products WHERE id = $1", id)
	product, err := scanRowIntoProduct(row)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *Store) GetAllProducts()(*[]Product, error) {
	rows, err := s.db.Query("SELECT id, name, description, price, stock_quanity, image, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}
	products, err := s.ScanRowsIntoProduct(rows)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Store) InsertEmbedding(id uuid.UUID, embedding []float32) error {
	embeddingJSON, err := json.Marshal(embedding)
	_, err = s.db.Exec("UPDATE products SET product_description_embedding = $2 WHERE id = $1", id, embeddingJSON)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoProduct(row *sql.Row) (*Product, error) {
	product := new(Product)
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return product, nil
}

func (s *Store) ScanRowsIntoProduct(rows *sql.Rows) (*[]Product, error) {
	
	products := make([]Product, 0)

	for rows.Next() {
		product := Product{}

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return &products, nil
}