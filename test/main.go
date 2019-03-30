package main

import (
	"encoding/json"
    "io/ioutil"
	"fmt"
	//"io"
	"net/http"
	"github.com/gorilla/mux"
	//"strconv"
	"time"
)

type TestData struct {
	A string
	B string
	C string
}

func getTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Прив"))
}

func postTestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RX")
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	fmt.Println(err)
	var t TestData
    err = json.Unmarshal(body, &t)
    if err != nil {
        panic(err)
    }
    fmt.Println(t.A)
    fmt.Println(t.B)
    fmt.Println(t.C)
	
	//fmt.Println(r.FormValue("a"))
	//fmt.Println(r.FormValue("b"))
	//fmt.Println(r.FormValue("c"))
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/test", getTestHandler).Methods("GET")
	r.HandleFunc("/api/test", postTestHandler).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at http://127.0.0.1:8000/")
	fmt.Println(srv.ListenAndServe())
}