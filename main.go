package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"./router"
)

func main() {

	// initialize and connect with DB

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