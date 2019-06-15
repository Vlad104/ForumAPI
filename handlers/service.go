package handlers

import (
	"net/http"
	"../database"
)


// /service/status Получение инфомарции о базе данных
func GetStatus(w http.ResponseWriter, r *http.Request) {

	result := database.GetStatusDB()
	resp, err := result.MarshalJSON()

	switch err {
	case nil:
		makeResponse(w, 200, resp)
	default:
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /service/clear Очистка всех данных в базе
func Clear(w http.ResponseWriter, r *http.Request) {
	database.ClearDB()
	makeResponse(w, 200, []byte("Очистка базы успешно завершена"))
}

