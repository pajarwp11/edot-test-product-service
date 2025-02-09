package stock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

func (s *StockHttpRepository) GetAvailableStock(productList *[]product.GetAvailableStock) (map[int]int, error) {
	url := "http://localhost:8003/product-warehouse/available-stock"

	jsonData, err := json.Marshal(productList)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed get available stock")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	formattedData, ok := response.Data.(map[int]int)
	if !ok {
		return nil, errors.New("wrong response format")
	}
	return formattedData, nil
}
