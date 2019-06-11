package handlers

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../database"
	"../models"
	"github.com/gorilla/mux"
	"github.com/go-openapi/swag"
)

// /forum/create Создание форума
func CreateForum(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		return
	}	
	forum := &models.Forum{}
	err = json.Unmarshal(body, &forum)

	// reg := strfmt.NewFormats()
	// err = forum.Validate(reg)
	if err != nil {
		return
	}

	result, err := database.CreateForumDB(forum)
	resp, _ := result.MarshalBinary()

	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(forum.User)))
	case database.ForumIsExist:
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /forum/{slug}/details Получение информации о форуме
func GetForum(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	result, err := database.GetForumDB(slug)

	resp, _ := result.MarshalBinary()
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte(makeErrorForum(slug)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /forum/{slug}/create Создание ветки
func CreateForumThread(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		return
	}	
	thread := &models.Thread{}
	err = json.Unmarshal(body, &thread)
	thread.Forum = slug // иначе не знаю как

	//err = forum.Validate()
	if err != nil {
		return
	}

	result, err := database.CreateForumThreadDB(thread)

	resp, _ := result.MarshalBinary()
	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.ForumOrAuthorNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(slug)))
	case database.ThreadIsExist:
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}


// /forum/{slug}/threads Список ветвей обсужления форума
func GetForumThreads(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	queryParams := r.URL.Query()
	var limit, since, desc string
	if limit = queryParams.Get("limit"); limit == "" {
		limit = "1";
	}
	if since = queryParams.Get("since"); limit == "" {
		since = "";
	}
	if desc = queryParams.Get("desc"); limit == ""{
		desc = "false";
	}

	result, err := database.GetForumThreadsDB(slug, limit, since, desc)
	
	// resp, _ := result.MarshalBinary()
	resp, _ := swag.WriteJSON(result)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte(makeErrorForum(slug)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /forum/{slug}/users Пользователи данного форума
func GetForumUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	queryParams := r.URL.Query()
	var limit, since, desc string
	if limit = queryParams.Get("limit"); limit == "" {
		limit = "1";
	}
	if since = queryParams.Get("since"); since == "" {
		since = "";
	}
	if desc = queryParams.Get("desc"); desc == ""{
		desc = "false";
	}

	result, err := database.GetForumUsersDB(slug, limit, since, desc)
	
	// resp, _ := result.MarshalBinary()
	resp, _ := swag.WriteJSON(result)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(slug)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}