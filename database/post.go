package database

import (
	"fmt"
	"strconv"
	"../models"
)

const (
	getPostSQL = `
		SELECT id, author, message, forum, thread, created, "isEdited" 
		FROM posts 
		WHERE id = $1
	`
	updatePostSQL = `
		UPDATE posts 
		SET message = COALESCE($2, message), "isEdited" = ($2 IS NOT NULL AND $2 <> message) 
		WHERE id = $1 
		RETURNING author::text, created, forum, "isEdited", thread, message
	`
)

// /post/{id}/details Получение информации о ветке обсуждения
func GetPostDB(id int) (*models.Post, error) {
	fmt.Println("GetPostDB")

	post := models.Post{}

	rows := DB.pool.QueryRow(getPostSQL, id)

	err := rows.Scan(
		&post.ID,
		&post.Author,
		&post.Message,
		&post.Forum,
		&post.Thread,
		&post.Created,
		&post.IsEdited,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, PostNotFound
		}
		return nil, err
	}

	return &post, nil
}


// /post/{id}/details Получение информации о ветке обсуждения
func GetPostFullDB(id int, related []string) (*models.PostFull, error) {
	fmt.Println("GetPostFullDB")
	postFull := models.PostFull{}
	var err error
	postFull.Post, err = GetPostDB(id)
	if err != nil {
		return nil, err
	}

	fmt.Println(related)
	for _, model := range related {
		switch model {
		case "thread":
			postFull.Thread, err = GetThreadDB(strconv.Itoa(int(postFull.Post.Thread)))
		case "forum":
			postFull.Forum, err = GetForumDB(postFull.Post.Forum)
		case "user":
			fmt.Println("user")
			postFull.Author, err = GetUserDB(postFull.Post.Author)
			fmt.Println(err)
		}

		if err != nil {
			return nil, err
		}
	}

	return &postFull, nil
}

// /post/{id}/details Изменение сообщения
func UpdatePostDB(postUpdate *models.PostUpdate, id int) (*models.Post, error) {
	fmt.Println("UpdatePostDB")
	post, err := GetPostDB(id)
	if err != nil {
		return nil, PostNotFound
	}

	if len(postUpdate.Message) == 0 {
		return post, nil
	}

	rows := DB.pool.QueryRow(updatePostSQL, strconv.Itoa(id), &postUpdate.Message)

	err = rows.Scan(
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.IsEdited,
		&post.Thread,
		&post.Message,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, PostNotFound
		}
		return nil, err
	}

	return post, nil
}
