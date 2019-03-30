package service

import (	
	"fmt"
	"net/http"
	"io/ioutil"
	"../database"
	"encoding/json"
	//"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func CreateForum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RX")
	body, err := ioutil.ReadAll(r.Body)	
	fmt.Println(string(body))
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	

	forum := &models.Forum{}
	err = json.Unmarshal(body, &forum)

	fmt.Println(forum.Slug)
	fmt.Println(forum.Title)
	fmt.Println(forum.User)

	//err = forum.Validate()

	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.CreateForumDB(forum)
	resp, err1 := result.MarshalBinary()
	fmt.Println(err1)
	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte("Hello1 "))
	case database.ForumIsExist:
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}
