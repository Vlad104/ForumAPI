package service

import (	
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../database"
	"../models"
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
		fmt.Println(err)
		return
	}	
	user := &models.User{}
	err = json.Unmarshal(body, &user)
	user.Nickname = nickname

	//err = forum.Validate()
	if err != nil {
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
		makeResponse(w, 500, []byte("Hello here"))
	}
}


// /user/{nickname}/profile Получение информации о пользователе
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nickname := params["nickname"]

	result, err := database.GetUserDB(nickname)
	fmt.Println(result)
	fmt.Println(err)

	resp, _ := result.MarshalBinary()

	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(nickname)))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}

// /user/{nickname}/profile Изменение данных о пользователе
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateUser")
	params := mux.Vars(r)
	nickname := params["nickname"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	user := &models.User{}
	err = json.Unmarshal(body, &user)
	user.Nickname = nickname

	//err = forum.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = database.UpdateUserDB(user)
	fmt.Println(err)
	resp, _ := user.MarshalBinary()

	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	case database.UserUpdateConflict:
		makeResponse(w, 409, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Can't find user with id #42\n"))
	}
}

