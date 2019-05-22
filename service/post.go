package service

import (
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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return
	}

	
	queryParams := r.URL.Query()
	relatedQuery := queryParams.Get("related")
	related := []string{}
	related = append(related, strings.Split(string(relatedQuery), ",")...)

	result, err := database.GetPostFullDB(id, related)

	resp, _ := result.MarshalBinary()
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
		return
    }

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		return
	}	
	postUpdate := &models.PostUpdate{}
	err = json.Unmarshal(body, &postUpdate)

	if err != nil {
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
