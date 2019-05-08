package database

import (
	// "bytes"
	"fmt"
	"strconv"
	"github.com/bozaro/tech-db-forum/generated/models"
	// "../models"
	// "time"
	// strfmt "github.com/go-openapi/strfmt"
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

	// var datetime time.Time
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
		&thread.Created,
		// &datetime,
	)
	// pdatetime := strfmt.DateTime(datetime)
	// thread.Created = &pdatetime
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
