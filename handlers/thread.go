package handlers

import (
	// "fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/Vlad104/TP_DB_RK2/database"
	"github.com/Vlad104/TP_DB_RK2/models"
	"github.com/gorilla/mux"
	"github.com/go-openapi/swag"
)

// /thread/{slug_or_id}/details Получение информации о ветке обсуждения
func GetThread(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["slug_or_id"]

	result, err := database.GetThreadDB(param)

	switch err {
	case nil:
		resp, _ := result.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.ThreadNotFound:
		makeResponse(w, 404, []byte(makeErrorThread(param)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/details Обновление ветки
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["slug_or_id"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}	
	threadUpdate := &models.ThreadUpdate{}
	err = threadUpdate.UnmarshalJSON(body)

	//err = forum.Validate()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}

	result, err := database.UpdateThreadDB(threadUpdate, param)

	switch err {
	case nil:
		resp, _ := result.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.PostNotFound:
		makeResponse(w, 404, []byte(makeErrorThread(param)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/create Создание новых постов
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("/thread/{slug_or_id}/create")
	params := mux.Vars(r)
	param := params["slug_or_id"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}	
	posts := &models.Posts{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}

	result, err := database.CreateThreadDB(posts, param)

	resp, _ := swag.WriteJSON(result)

	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.ThreadNotFound:
		makeResponse(w, 404, []byte(makeErrorThreadID(param)))
	case database.UserNotFound:
		makeResponse(w, 404, []byte(makeErrorPostAuthor(param)))
	case database.PostParentNotFound:
		makeResponse(w, 409, []byte(makeErrorThreadConflict()))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/posts Сообщения данной ветви обсуждения
func GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["slug_or_id"]
	queryParams := r.URL.Query()
	var limit, since, sort, desc string
	if limit = queryParams.Get("limit"); limit == "" {
		limit = "1";
	}
	since = queryParams.Get("since"); 
	// if since = queryParams.Get("since"); since == "" {
	// 	since = "";
	// }
	if sort = queryParams.Get("sort"); sort == ""{
		sort = "flat";
	}
	if desc = queryParams.Get("desc"); desc == ""{
		desc = "false";
	}
	// fmt.Println("limit", limit, "since", since, "sort", sort, "desc", desc)
	result, err := database.GetThreadPostsDB(param, limit, since, sort, desc)
	
	// resp, _ := result.MarshalJSON()
	resp, _ := swag.WriteJSON(result)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte(makeErrorThread(param)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/vote Проголосовать за ветвь обсуждения
func MakeThreadVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	param := params["slug_or_id"]
	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		makeResponse(w, 500, []byte(err.Error()))
		return
	}	
	vote := &models.Vote{}
	err = vote.UnmarshalJSON(body)

	result, err := database.MakeThreadVoteDB(vote, param)	
	
	switch err {
	case nil:
		resp, _ := result.MarshalJSON()
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte(makeErrorThread(param)))
	case database.UserNotFound:
		makeResponse(w, 404, []byte(makeErrorUser(param)))
	default:		
		makeResponse(w, 500, []byte(err.Error()))
	}
}
