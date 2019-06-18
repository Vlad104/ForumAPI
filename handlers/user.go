package handlers

import (
	"net/http"
	"io/ioutil"
	"github.com/Vlad104/TP_DB_RK2/database"
	"github.com/Vlad104/TP_DB_RK2/models"
	"github.com/gorilla/mux"
	"github.com/go-openapi/swag"
)

// /user/{nickname}/create Создание нового пользователя
func CreateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nickname := params["nickname"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}	
	user := &models.User{}
	err = user.UnmarshalJSON(body)
	user.Nickname = nickname

	//err = forum.Validate()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}
	result, err := database.CreateUserDB(user)

	switch err {
	case nil:
		resp, _ := swag.WriteJSON(user)
		makeResponse(w, 201, resp)
	case database.UserIsExist:
		resp, _ := swag.WriteJSON(result)
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /user/{nickname}/profile Получение информации о пользователе
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nickname := params["nickname"]

	result, err := database.GetUserDB(nickname)

	switch err {
	case nil:
		resp, _ := result.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(nickname)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /user/{nickname}/profile Изменение данных о пользователе
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nickname := params["nickname"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}	
	user := &models.User{}
	err = user.UnmarshalJSON(body)
	user.Nickname = nickname

	//err = forum.Validate()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}
	err = database.UpdateUserDB(user)

	switch err {
	case nil:
		resp, _ := user.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(nickname)))
	case database.UserUpdateConflict:
		makeResponse(w, 409, []byte(makeErrorEmail(nickname)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

