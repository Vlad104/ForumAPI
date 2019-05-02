package database

import (
	// "bytes"
	"github.com/bozaro/tech-db-forum/generated/models"
)

const (
	getThreadSQL = `
		SELECT id, title, author, forum, message, votes, slug, created
		FROM threads
		WHERE slug = $1
	`
)

func GetThreadDB(slug string) (*models.Thread, error) {
	var err error
	var thread models.Thread

	err = DB.pool.QueryRow(
		getThreadSQL, 
		slug,
	).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	if err != nil {
		return nil, ThreadNotFound
	}

	return &thread, nil
}