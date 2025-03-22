package llmclient

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"math/rand"

	"github.com/kidusshun/ecom_bot/config"
	"github.com/kidusshun/ecom_bot/embedding"
	"github.com/kidusshun/ecom_bot/service/product"
	"github.com/pgvector/pgvector-go"
)

type QueryStore struct {
	DB *sql.DB
	ProductStore product.Store
}

func NewQueryStore(db *sql.DB) *QueryStore {
	return &QueryStore{
		DB: db,
	}
}

func (s *QueryStore) CompanyInfo(query string) (*ToolCallResponse, error) {
	queryStr := fmt.Sprintf("SELECT * FROM documents ORDER BY embedding <=> $1 LIMIT $2")
	embedding, err := embedding.GetEmbedding(query)
	if err != nil {
		return nil, err
	}

	vector := pgvector.NewVector(embedding.Embedding)

	// Execute query
	rows, err := s.DB.Query(queryStr, vector, 5)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Process rows and return the results
	var results []map[string]string
	for rows.Next() {
		var id int
		var text string
		var embedding pgvector.Vector // Assuming pgvector is used to handle the vector column
		if err := rows.Scan(&id, &text, &embedding); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		results = append(results, map[string]string{"text": text})
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
	modelResponse, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	response := Message{
		Role: USER,
		Parts: []Part{
			{
				FunctionResponse: &FunctionResponse{
					Name: "QueryProducts",
					Response: Result{
						Result: string(modelResponse),
					},
				},
			},
		},
	}
	return &ToolCallResponse{ModelResponse: response}, nil	
}


func (s *QueryStore) QueryProducts(query string) (*ToolCallResponse, error)  {
	queryStr := fmt.Sprintf("SELECT id, name, description, price, stock_quanity, image, created_at, updated_at FROM products ORDER BY product_description_embedding <=> $1 LIMIT $2")
	embedding, err := embedding.GetEmbedding(query)
	if err != nil {
		return nil, err
	}

	vector := pgvector.NewVector(embedding.Embedding)

	// Execute query
	rows, err := s.DB.Query(queryStr, vector, 5)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products, err:= s.ProductStore.ScanRowsIntoProduct(rows)

	productsNoImages := make([]product.ProductNOImage, 0)
	for _, producte := range *products {
		productImage := product.ProductNOImage{
			ID: producte.ID,
			Name: producte.Name,
			Description: producte.Description,
			Price: producte.Price,
			StockQuantity: producte.StockQuantity,
			CreatedAt: producte.CreatedAt,
			UpdatedAt: producte.UpdatedAt,
		}
		productsNoImages = append(productsNoImages, productImage)
	}

	log.Println(len(productsNoImages), " products gotten")

	if err != nil {
		return nil, err
	}
	
	modelResponse, err := json.Marshal(productsNoImages)
	if err != nil {
		return nil, err
	}
	response := Message{
		Role: USER,
		Parts: []Part{
			{
				FunctionResponse: &FunctionResponse{
					Name: "QueryProducts",
					Response: Result{
						Result: string(modelResponse),
					},
				},
			},
		},
	}
	return &ToolCallResponse{ModelResponse: response,Products: *products}, nil
}


func (s *QueryStore) TrackOrder(orderID string) (*ToolCallResponse, error)  {
	apiURL := "https://api.goshippo.com/tracks/"
	statuses := []string{"SHIPPO_RETURNED", "SHIPPO_PRE_TRANSIT", "SHIPPO_DELIVERED", "SHIPPO_RETURNED"}
	randomIndex := rand.Intn(len(statuses))

	formData := url.Values{}
	formData.Set("carrier", "shippo")
	formData.Set("tracking_number", statuses[randomIndex]) // Random status
	formData.Set("metadata", "Test shipment")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// Set headers correctly
	req.Header.Set("Authorization", config.ShippoEnvs.ShippoAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(bodyBytes)) 
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var jsonObject map[string]interface{}
	err = json.Unmarshal(body, &jsonObject)

	jsonString, err := json.Marshal(jsonObject)
	if err != nil {
		fmt.Println("Error converting to JSON string:", err)
		return nil, err
	}
	response := Message{
		Role: USER,
		Parts: []Part{
			{
				FunctionResponse: &FunctionResponse{
					Name: "TrackOrder",
					Response: Result{
						Result: string(jsonString),
					},
				},
			},
		},
	}
	return &ToolCallResponse{
		ModelResponse: response,
	}, nil

}