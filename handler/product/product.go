package product

import (
	"encoding/json"
	"net/http"
	"product-service/models/product"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type ProductUsecase interface {
	Register(productRegister *product.RegisterRequest) error
	GetList(request *product.GetListRequest) (*[]product.Product, int, error)
}

type ProductHandler struct {
	productUsecase ProductUsecase
}

type Meta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta,omitempty"`
}

var validate = validator.New()

func NewProductHandler(productUsecase ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
	}
}

func (p *ProductHandler) Register(w http.ResponseWriter, req *http.Request) {
	request := product.RegisterRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := p.productUsecase.Register(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "product registered"
	json.NewEncoder(w).Encode(response)
}

func (p *ProductHandler) GetList(w http.ResponseWriter, req *http.Request) {
	request := product.GetListRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")

	request.Category = req.URL.Query().Get("category")
	request.ShopId, _ = strconv.Atoi(req.URL.Query().Get("shop_id"))
	request.Page, _ = strconv.Atoi(req.URL.Query().Get("page"))
	request.PerPage, _ = strconv.Atoi(req.URL.Query().Get("per_page"))

	if request.Page == 0 {
		request.Page = 1
	}

	if request.PerPage == 0 {
		request.PerPage = 10
	}

	products, total, err := p.productUsecase.GetList(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Message = "get product list success"
	response.Data = map[string]interface{}{
		"products": products,
	}
	response.Meta = &Meta{
		Page:      request.Page,
		PerPage:   request.PerPage,
		Total:     total,
		TotalPage: ((total + request.PerPage - 1) / request.PerPage),
	}
	json.NewEncoder(w).Encode(response)
}
