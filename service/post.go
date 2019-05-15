package service

import (	
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../database"
	"../models"
	"github.com/gorilla/mux"
)

// /post/{id}/details Получение информации о ветке обсуждения
func GetPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetPost")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
    if err != nil {
		fmt.Println(err)
		return
    }

	result, err := database.GetPostDB(id)

	resp, _ := result.MarshalBinary()
	fmt.Println("DB result")
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}

// /post/{id}/details Изменение сообщения
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdatePost")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
    if err != nil {
		fmt.Println(err)
		return
    }

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	postUpdate := &models.PostUpdate{}
	err = json.Unmarshal(body, &postUpdate)

	//err = forum.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.UpdatePostDB(postUpdate, id)

	resp, _ := result.MarshalBinary()
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}
