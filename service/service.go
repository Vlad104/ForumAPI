package service

import (	
	"fmt"
	// "strconv"
	"net/http"
	// "io/ioutil"
	"../database"
	// "encoding/json"
	//"github.com/bozaro/tech-db-forum/generated/client/operations"
	// "github.com/bozaro/tech-db-forum/generated/models"
	// "github.com/gorilla/mux"
	// "github.com/go-openapi/swag"
)


// /service/status Получение инфомарции о базе данных
func GetStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetStatus")

	result := database.GetStatusDB()
	resp, err := result.MarshalBinary()

	switch err {
	case nil:
		makeResponse(w, 200, resp)
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}


// /service/clear Очистка всех данных в базе
func Clear(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Clear")
	database.ClearDB()
	makeResponse(w, 200, []byte("Очистка базы успешно завершена"))
}

