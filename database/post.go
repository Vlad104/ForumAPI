package database

import (
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

// /post/{id}/details Изменение сообщения
func UpdatePostDB(postUpdate *models.PostUpdate, id int) (*models.Post, error) {
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
