package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./database"
	"./handlers"
)

func main() {

	database.DB.Connect()

	router := handlers.CreateRouter()	

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:5000",
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}
	fmt.Println("Starting server at PORT: 5000")
	log.Fatal(srv.ListenAndServe())
}
