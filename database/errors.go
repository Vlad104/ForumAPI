package database

import (
	"encoding/json"
	"log"
)

const (
	pgxOK				= ""
	pgxErrNotNull		= "23502"
	pgxErrForeignKey 	= "23503"
	pgxErrUnique 		= "23505"
)

type Error struct {
	Message string `json:"message,omitempty"`
}

func (r *Error) Error() string {
	errorBytes, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	return string(errorBytes)
}

func New(msg string) error {
	return &Error{Message: msg}
}

// ForumIsExist error.
var ForumIsExist = New("Forum was created earlier")

// ForumNotFound error.
var ForumNotFound = New("Forum not found")

// ForumOrAuthorNotFound error.
var ForumOrAuthorNotFound = New("Forum or Author not found")

// UserNotFound error.
var UserNotFound = New("User not found")

// UserIsExist error.
var UserIsExist = New("User was created earlier")

// UserUpdateConflict error.
var UserUpdateConflict = New("User not updated")

// ThreadIsExist error.
var ThreadIsExist = New("Thread was created earlier")

// ThreadNotFound error.
var ThreadNotFound = New("Thread not found")

// PostParentNotFound error.
var PostParentNotFound = New("No parent for thread")

// PostNotFound error.
var PostNotFound = New("Post not found")
