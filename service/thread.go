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
		makeResponse(w, 404, []byte(makeErrorThread(param)))
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
		makeResponse(w, 404, []byte(makeErrorThread(param)))
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

	result, err := database.CreateThreadDB(posts, param)

	resp, _ := swag.WriteJSON(result)
	fmt.Println(err)
	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.ThreadNotFound:
		makeResponse(w, 404, []byte(makeErrorThreadID(param)))
	case database.UserNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	case database.PostNotFound:
		makeResponse(w, 409, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}

// НЕ ТЕСТИРОВАЛ
// /thread/{slug_or_id}/posts Сообщения данной ветви обсуждения
func GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetThreadPosts")
	params := mux.Vars(r)
	param := params["slug_or_id"]
	queryParams := r.URL.Query()
	var limit, since, sort, desc string
	if limit = queryParams.Get("limit"); limit == "" {
		limit = "1";
	}
	if since = queryParams.Get("since"); since == "" {
		since = "";
	}
	if sort = queryParams.Get("sort"); sort == ""{
		sort = "flat";
	}
	if desc = queryParams.Get("desc"); desc == ""{
		desc = "false";
	}
	fmt.Println(param, limit, since, sort, desc)

	result, err := database.GetThreadPostsDB(param, limit, since, sort, desc)
	fmt.Println(result)
	fmt.Println(err)
	
	// resp, _ := result.MarshalBinary()
	resp, _ := swag.WriteJSON(result)
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte(makeErrorThread(param)))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}


// НЕ ТЕСТИРОВАЛ
// /thread/{slug_or_id}/vote Проголосовать за ветвь обсуждения
func MakeThreadVote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MakeThreadVoteDB")
	params := mux.Vars(r)
	param := params["slug_or_id"]
	
	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	vote := &models.Vote{}
	err = json.Unmarshal(body, &vote)

	result, err := database.MakeThreadVoteDB(vote, param)
	fmt.Println(result)
	fmt.Println(err)
	
	resp, _ := result.MarshalBinary()
	// resp, _ := swag.WriteJSON(result)
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err)
	
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}
