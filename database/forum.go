package database

import (
	//"bytes"
	"github.com/jackc/pgx"	
	"github.com/bozaro/tech-db-forum/generated/models"
)

const (
	createForumSQL = `
		INSERT
		INTO forums (slug, title, "user")
		VALUES ($1, $2, (SELECT nickname FROM users WHERE nickname = $3)) 
		RETURNING "user"`
	getForumSQL = `
		SELECT slug, title, "user", 
			(SELECT COUNT(*) FROM posts WHERE forum = $1), 
			(SELECT COUNT(*) FROM threads WHERE forum = $1)
		FROM forums
		WHERE slug = $1`
)

func CreateForumDB(f *models.Forum) (*models.Forum, error)  {
	rows := DB.pool.QueryRow(
		createForumSQL,
		&f.Slug,
		&f.Title,
		&f.User,
	)

	err := rows.Scan(&f.User)
	if err != nil {
		switch err.(pgx.PgError).Code {
		case pgxErrUnique:
			forum, _ := GetForumDB(f.Slug)
			return forum, ForumIsExist
		case pgxErrNotNull:
			return nil, UserNotFound
		default:
			return nil, err
		}
	}
	return f, nil
}

func GetForumDB(slug string) (*models.Forum, error) {
	forum := models.Forum{}

	err := DB.pool.QueryRow(getForumSQL, slug).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
	)

	if err != nil {
		return nil, ForumNotFound
	}

	return &forum, nil
}