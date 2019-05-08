package router

import (
	"../service"
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/forum/create", service.CreateForum).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/create", service.CreateForumThread).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", service.GetForum).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", service.GetForumThreads).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", service.GetForumUsers).Methods("GET")
	r.HandleFunc("/api/user/{nickname}/create", service.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", service.GetUser).Methods("GET")
	r.HandleFunc("/api/user/{nickname}/profile", service.UpdateUser).Methods("POST")
	r.HandleFunc("/api/post/{id}/details", service.GetPost).Methods("GET")
	r.HandleFunc("/api/post/{id}/details", service.UpdatePost).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", service.GetThread).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/details", service.UpdateThread).Methods("POST")
	r.HandleFunc("/api/service/status", service.GetStatus).Methods("GET")
	r.HandleFunc("/api/service/clear", service.Clear).Methods("POST")
/*
	r.HandleFunc("/api/thread/{slug_or_id}/create", service.createThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", service.getPosts).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", service.postVote).Methods("POST")
*/
	return r
}
