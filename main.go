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
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./common/"))))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./common/swagger-ui/"))))

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at http://127.0.0.1:5000")
	log.Fatal(srv.ListenAndServe())
}
