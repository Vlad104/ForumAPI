package router

import (
	"../service"
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/forum/create", service.CreateForum).Methods("POST")
/*
	r.HandleFunc("/api/forum/{slug}/createBranch", createBranch).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", getDetails).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", getThreads).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", getUsers).Methods("GET")
	r.HandleFunc("/api/post/{id}/details", getDetails).Methods("GET")
	r.HandleFunc("/api/post/{id}/details", postDetails).Methods("POST")
	r.HandleFunc("/api/service/clear", clearService).Methods("POST")
	r.HandleFunc("/api/service/status", getStatus).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/create", createThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", getThread).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/details", updateThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", getPosts).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", postVote).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/create", createUser).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", getUser).Methods("GET")
	r.HandleFunc("/api/user/{nickname}/profile", updateUser).Methods("GET")
*/
	return r
}
