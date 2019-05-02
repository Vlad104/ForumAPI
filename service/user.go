package service

import (	
	"fmt"
	"net/http"
	"io/ioutil"
	"../database"
	"encoding/json"
	//"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/gorilla/mux"
	"github.com/go-openapi/swag"
)

// /user/{nickname}/create Создание нового пользователя
func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateUser")
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
	result, err := database.CreateUserDB(user)
	fmt.Println(result)
	fmt.Println(err)

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
	fmt.Println("GetUser")
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
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
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
	case database.UserIsExist:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	case database.UserIsExist:
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte("Can't find user with id #42\n"))
	}
}

