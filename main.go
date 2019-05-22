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
		Addr:         "0.0.0.0:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at PORT: 5000")
	log.Fatal(srv.ListenAndServe())
}
