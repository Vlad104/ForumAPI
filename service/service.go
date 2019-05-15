package service

import (	
	"fmt"
	"net/http"
	"../database"
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

