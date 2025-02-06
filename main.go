package main

import (
	"fmt"
	"log"
	"net/http"
	"product-service/db/mysql"

	"github.com/gorilla/mux"
)

func main() {
	mysql.Connect()
	router := mux.NewRouter()

	fmt.Println("server is running")
	err := http.ListenAndServe(":8001", router)
	if err != nil {
		log.Fatal(err)
	}
}
