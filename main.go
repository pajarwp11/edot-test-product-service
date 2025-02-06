package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	fmt.Println("server is running")
	err := http.ListenAndServe(":8001", router)
	if err != nil {
		log.Fatal(err)
	}
}
