package stock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"product-service/models/product"
)

type StockHttpRepository struct {
}

func NewStockHttpRepository() *StockHttpRepository {
	return &StockHttpRepository{}
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *StockHttpRepository) GetAvailableStock(productList *[]product.GetAvailableStock) (map[string]interface{}, error) {
	url := "http://localhost:8003/product-warehouse/available-stock"

	jsonData, err := json.Marshal(productList)
	if err != nil {
		log.Printf("failed to marshal request: %v\n", err)
		return nil, err
	}

	log.Println("request body:", string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("error init request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error request available stock: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v\n", err)
		return nil, err
	}

	log.Println("response body:", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Println("failed to get available stock, status code:", resp.StatusCode)
		return nil, errors.New("failed get available stock")
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("failed to unmarshal response: %v\n", err)
		return nil, err
	}
	formattedData, _ := response.Data.(map[string]interface{})
	return formattedData, nil
}
