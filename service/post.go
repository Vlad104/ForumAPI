package service

import (	
	"fmt"
	"strings"
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
	id, err := strconv.ParseInt(params["id"], 10, 64)
	fmt.Println(params["id"])
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	
	fmt.Println(r.URL)
	queryParams := r.URL.Query()
	fmt.Println(queryParams)
	relatedQuery := queryParams.Get("related")
	fmt.Println(relatedQuery)
	// related := []string{"post"}
	related := []string{}
	related = append(related, strings.Split(string(relatedQuery), ",")...)
	fmt.Println(related)

	// result_temp, err := database.GetPostDB(id, related)
	// result := models.PostFull{}
	// result.Post = result_temp

	result, err := database.GetPostFullDB(id, related)

	resp, _ := result.MarshalBinary()
	fmt.Println("DB result")
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte(makeErrorPost(string(id))))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}

// /post/{id}/details Изменение сообщения
func UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := database.UpdatePostDB(postUpdate, id)
	resp, _ := result.MarshalBinary()
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte(makeErrorPost(string(id))))
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}
