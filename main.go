package main

import (
	"fmt"
	"log"
	"net/http"
	"product-service/db/mysql"
	productHandler "product-service/handler/product"
	"product-service/middleware"
	"product-service/repository/http/stock"
	productRepo "product-service/repository/product"
	productUsecase "product-service/usecase/product"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()

	router := mux.NewRouter()
	stockHttpRepo := stock.NewStockHttpRepository()
	productRepository := productRepo.NewProductRepository(mysql.MySQL)
	productUsecase := productUsecase.NewProductUsecase(productRepository, stockHttpRepo)
	productHandler := productHandler.NewProductHandler(productUsecase)
	router.Handle("/product/register", middleware.JWTMiddleware(http.HandlerFunc(productHandler.Register))).Methods(http.MethodPost)
	router.Handle("/product", middleware.JWTMiddleware(http.HandlerFunc(productHandler.Register))).Methods(http.MethodGet)

	fmt.Println("server is running")
	err := http.ListenAndServe(":8001", router)
	if err != nil {
		log.Fatal(err)
	}
}
