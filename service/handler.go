package service

import (
	"fmt"
	"net/http"
)

func makeResponse(w http.ResponseWriter, status int, resp []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func makeErrorUser(s string) string {
	return fmt.Sprintf(`{"message": "Can't find user by nickname: %s"}`, s)
}

func makeErrorForum(s string) string {
	return fmt.Sprintf(`{"message": "Can't find forum with slug: %s"}`, s)
}

func makeErrorThread(s string) string {
	return fmt.Sprintf(`{"message": "Can't find thread by slug: %s"}`, s)
}