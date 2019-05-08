package database

import (
	// "bytes"
	"fmt"
	"strconv"
	"github.com/bozaro/tech-db-forum/generated/models"
	// "../models"
	"time"
	//strfmt "github.com/go-openapi/strfmt"
	strmft "github.com/bozaro/tech-db-forum/vendor/github.com/go-openapi/strfmt"
)

const (
	getThreadSlugSQL = `
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE slug = $1
	`
	getThreadIdSQL = `
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE id = $1
	`
	updateThreadSQL = `
		UPDATE threads
		SET title = coalesce(nullif($2, ''), title),
			message = coalesce(nullif($3, ''), message)
		WHERE slug = $1
		RETURNING id, title, author, forum, message, votes, slug, created
	`
	// createPostSQL = `
	// 	INSERT 
	// 	INTO posts (author, created, message, thread, parent, forum, path) 
	// 	VALUES ($1, $2, $3, $4, $5, $6, 
	// 		(SELECT path FROM posts WHERE id = %d) || 
	// 		(select currval(pg_get_serial_sequence('posts', 'id')))
	// 	)
	// 	RETURNING author, created, forum, id, message, parent, thread
	// `
	// createPostDelSQL = `
	// 	INSERT 
	// 	INTO posts (author, created, message, thread, parent, forum, path) 
	// 	VALUES ($1, $2, $3, $4, $5, $6, 
	// 		(SELECT path FROM posts WHERE id = %d) || 
	// 		(select currval(pg_get_serial_sequence('posts', 'id')))
	// 	),
	// 	RETURNING author, created, forum, id, message, parent, thread
	// `
)

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func GetThreadDB(param string) (*models.Thread, error) {
	var err error
	var thread models.Thread

	query := getThreadSlugSQL
	if isNumber(param) {
		query = getThreadIdSQL
	}

	var datetime time.Time
	err = DB.pool.QueryRow(
		query,
		param,
	).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		// &thread.Created,
		&datetime,
	)
	pdatetime := strfmt.DateTime(datetime)
	thread.Created = &pdatetime
	fmt.Println(thread)
	fmt.Println(err)
	if err != nil {
		return nil, ThreadNotFound
	}

	return &thread, nil
}

// /thread/{slug_or_id}/details Обновление ветки
func UpdateThreadDB(thread *models.ThreadUpdate, param string) (*models.Thread, error) {
	threadFound, err := GetThreadDB(param)
	if err != nil {
		return nil, err
	}

	updatedThread := models.Thread{}

	err = DB.pool.QueryRow(updateThreadSQL,
		&threadFound.Slug,
		&thread.Title,
		&thread.Message,
	).Scan(
		&updatedThread.ID,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)

	if err != nil {
		return nil, err
	}

	return &updatedThread, nil
}

func authorExists(nickname string) bool {
	var user models.User
	rows :=  DB.pool.QueryRow(getUserByNickname, nickname)

	if err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if err.Error() == "no rows in result set" {
			return true
		}
	}
	return false
}

// переделать
func parentExitsInOtherThread(parent int64, threadID int) bool {
	var t int
	rows := DB.pool.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)`,
		parent, threadID)

	if err := rows.Scan(&t); err != nil {
		if err.Error() == "no rows in result set" {
			return false
		}
	}
	return true
}

// /thread/{slug_or_id}/create Создание новых постов
func CreateThreadDB(posts *models.Posts, param string) (*models.Posts, error) {
	thread, err := GetThreadDB(param)
	if err != nil {
		return nil, err
	}

	if len(*posts) == 0 {
		return posts, nil
	}

	// надо подумать
	// пока такой костыль
	created := time.Now().Format("2006-01-02 15:04:05")
	var query := strings.Builder{}
	query.WriteString("INSERT INTO posts (author, created, message, thread, parent, forum, path) VALUES ")
	queryBody := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (select currval(pg_get_serial_sequence('posts', 'id')))),"
	queryBodyEnd := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (select currval(pg_get_serial_sequence('posts', 'id'))))"
	for i, post := range *posts {
		if authorExists(post.Author) {
			return nil, UserNotFound
		}
		if parentExitsInOtherThread(post.Parent, thread.Id) || parentNotExists(post.Parent) {
			return nil, PostParentNotFound
		}

		// можно оптимизировать
		if i < len(*posts) - 1 {
			queryBuilder.WriteString(fmt.Sprintf(queryBody, post.Author, created, post.Message, thread.Id, post.Parent, thread.Forum, post.Parent))
		} else {
			queryBuilder.WriteString(fmt.Sprintf(queryBodyEnd, post.Author, created, post.Message, thread.Id, post.Parent, thread.Forum, post.Parent))
		}

	}
	query.WriteString("RETURNING author, created, forum, id, message, parent, thread")

	rows, err := DB.pool.Query(query.String()) 
	if err != nil {
		return nil, err
	}
	insertPosts := := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		_ = rows.Scan(
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.Id,
			&post.Message,
			&post.Parent,
			&post.Thread,
		)
		insertPosts = append(insertPosts, &post) 
	}
	return &insertPosts, nil
}