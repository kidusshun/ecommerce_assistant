package script

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/kidusshun/ecom_bot/embedding"
	"github.com/kidusshun/ecom_bot/service/product"
)


type ProductObj struct {
	Name  string `json:"product_name"`
	Price float64 `json:"price"`
	StockQuantity int `json:"stock_quantity"`
	ProductDescription string `json:"product_description"`
	Image string `json:"image"`
}

type Store struct {
	prodStore *product.Store
}

func NewStore(prodStore *product.Store) *Store {
	return &Store{
		prodStore: prodStore,
	}
}

func (s *Store) LoadProducts()([]ProductObj,error) {
	file, err := os.Open("/home/kidus/code/go/ai_apps/ecom_cutomer_bot/products.json")
	if err != nil {
		return nil, err
	}
	
	defer file.Close()
	
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	
	var products []ProductObj
	if err := json.Unmarshal(byteValue, &products); err != nil {
		return nil, err
	}

	for _, prod := range products {
		product := product.Product{
			Name: prod.Name,
			Description: prod.ProductDescription,
			Price: prod.Price,
			StockQuantity: prod.StockQuantity,
			Image: prod.Image,
		}
		err := s.prodStore.AddProduct(product)
		if err != nil {
			log.Println(err)
		}
	}

	return []ProductObj{}, nil
}


func (s *Store) InsertEmbedding() error {
	prods,err := s.prodStore.GetAllProducts()
	if err != nil {
		return err
	}

	for _, prod := range *prods {
		embedding, err := embedding.GetEmbedding(prod.Description)
		if err != nil {
			return err
		}
		err = s.prodStore.InsertEmbedding(prod.ID, embedding.Embedding)
		if err != nil {
			return err
		}
	}	
	
	return nil
}