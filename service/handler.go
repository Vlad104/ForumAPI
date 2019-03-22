package service

import (	
	"fmt"
	"net/http"
	"io/ioutil"
	"../database"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
)

func CreateForum(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	forum := operations.NewForumCreateParams()
	err = forum.Forum.UnmarshalBinary(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	//err = forum.Validate()

	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.CreateForumDB(forum.Forum)
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
