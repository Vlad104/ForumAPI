package database

import (
	"strconv"
	"github.com/Vlad104/TP_DB_RK2/models"
)

const (
	getPostSQL = `
		SELECT id, author, message, forum, thread, created, "isEdited", parent
		FROM posts 
		WHERE id = $1
	`
	updatePostSQL = `
		UPDATE posts 
		SET message = COALESCE($2, message), "isEdited" = ($2 IS NOT NULL AND $2 <> message) 
		WHERE id = $1 
		RETURNING author::text, created, forum, "isEdited", thread, message, parent
	`
)

// /post/{id}/details Получение информации о ветке обсуждения
func GetPostDB(id int) (*models.Post, error) {
	post := models.Post{}

	err := DB.pool.QueryRow(
		getPostSQL,
		id,
	).Scan(
		&post.ID,
		&post.Author,
		&post.Message,
		&post.Forum,
		&post.Thread,
		&post.Created,
		&post.IsEdited,
		&post.Parent,
	)

	if err == nil {
		return &post, nil
	} else if (err.Error() == noRowsInResult) {
		return nil, PostNotFound
	} else {
		return nil, err
	}
}


// /post/{id}/details Получение информации о ветке обсуждения
func GetPostFullDB(id int, related []string) (*models.PostFull, error) {
	postFull := models.PostFull{}
	var err error
	postFull.Post, err = GetPostDB(id)
	if err != nil {
		return nil, err
	}

	for _, model := range related {
		switch model {
		case "thread":
			postFull.Thread, err = GetThreadDB(strconv.Itoa(int(postFull.Post.Thread)))
		case "forum":
			postFull.Forum, err = GetForumDB(postFull.Post.Forum)
		case "user":
			postFull.Author, err = GetUserDB(postFull.Post.Author)
		}

		if err != nil {
			return nil, err
		}
	}

	return &postFull, nil
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
		&post.Parent,
	)

	if err == nil {
		return post, nil
	} else if (err.Error() == noRowsInResult) {
		return nil, PostNotFound
	} else {
		return nil, err
	}
}
