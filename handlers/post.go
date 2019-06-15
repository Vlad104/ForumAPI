package handlers

import (
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"../database"
	"../models"
	"github.com/gorilla/mux"
)

// /post/{id}/details Получение информации о ветке обсуждения
func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}
	
	queryParams := r.URL.Query()
	relatedQuery := queryParams.Get("related")
	related := []string{}
	related = append(related, strings.Split(string(relatedQuery), ",")...)

	result, err := database.GetPostFullDB(id, related)

	switch err {
	case nil:
		resp, _ := result.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte(makeErrorPost(string(id))))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /post/{id}/details Изменение сообщения
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
    if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
    }

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}	
	postUpdate := &models.PostUpdate{}
	err = postUpdate.UnmarshalJSON(body)

	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}
	result, err := database.UpdatePostDB(postUpdate, id)
	switch err {
	case nil:
		resp, _ := result.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte(makeErrorPost(string(id))))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}
