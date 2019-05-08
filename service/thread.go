package service

import (	
	"fmt"
	"net/http"
	"io/ioutil"
	"../database"
	"encoding/json"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/gorilla/mux"
	// "github.com/go-openapi/swag"
)

// /thread/{slug_or_id}/details Получение информации о ветке обсуждения
func GetThread(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetThread")
	params := mux.Vars(r)
	param := params["slug_or_id"]

	result, err := database.GetThreadDB(param)

	resp, _ := result.MarshalBinary()

	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ThreadNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}

// /thread/{slug_or_id}/details Обновление ветки
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateThread")
	params := mux.Vars(r)
	param := params["slug_or_id"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	threadUpdate := &models.ThreadUpdate{}
	err = json.Unmarshal(body, &threadUpdate)

	//err = forum.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.UpdateThreadDB(threadUpdate, param)

	resp, _ := result.MarshalBinary()
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}

// /thread/{slug_or_id}/create Создание новых постов
func CreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreatePost")
	params := mux.Vars(r)
	param := params["slug_or_id"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	posts := &models.Posts{}
	err = json.Unmarshal(body, &posts)

	//err = forum.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.CreatePostDB(posts, param)

	resp, _ := result.MarshalBinary()
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	case database.PostNotFound:
		makeResponse(w, 409, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}
