package service

import (
	"net/http"
)

func makeResponse(w http.ResponseWriter, status int, resp []byte) {
	w.Write(resp)
}
