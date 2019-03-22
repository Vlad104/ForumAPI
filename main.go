package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./database"
	"./router"
)

func main() {

	database.DB.Connect()

	router := router.CreateRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at http://127.0.0.1:5000")
	log.Fatal(srv.ListenAndServe())
}
