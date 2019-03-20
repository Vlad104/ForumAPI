package main


func createForum(w http.ResponseWritter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Error. Method")
		return
	}
	// var forum Forum
	// and start create forum
}

	r.HandleFunc("/api/forum/create", createForum).Methods("POST")
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