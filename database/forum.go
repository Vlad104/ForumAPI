package database

import (
	"bytes"
	"github.com/jackc/pgx"	
	"github.com/bozaro/tech-db-forum/generated/models"
)

const (
	createForumSQL = `
		INSERT
		INTO forums (slug, title, "user")
		VALUES ($1, $2, (
			SELECT nickname FROM users WHERE nickname = $3
		)) 
		RETURNING "user"
	`
	getForumSQL = `
		SELECT slug, title, "user", 
			(SELECT COUNT(*) FROM posts WHERE forum = $1), 
			(SELECT COUNT(*) FROM threads WHERE forum = $1)
		FROM forums
		WHERE slug = $1
	`
	createForumThreadSQL = `
		INSERT
		INTO threads (authot, created, forum, message, slug, title)
		VALUES($1, $2, $3, $4, $5, (
			SELECT slug FROM forums WHERE slug = $6
		)) 
		RETURNING author, created, forum, id, message, title
	`
	getForumThreadsSinceSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $3 AND created $1 $4::TEXT::TIMESTAMPTZ
		ORDER BY created $2
		LIMIT $5::TEXT::INTEGER
	`
	getForumThreadsSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $2
		ORDER BY created $1
		LIMIT $3::TEXT::INTEGER
	`
	getForumUsersSienceSQl = `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT author FROM threads WHERE forum = $3
				UNION
				SELECT author FROM posts WHERE forum = $3
			) 
			AND LOWER(nickname) $1 LOWER($2::TEXT)
		ORDER BY nickname $2
		LIMIT $5::TEXT::INTEGER
	`
	getForumUsersSQl = `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT author FROM threads WHERE forum = $2
				UNION
				SELECT author FROM posts WHERE forum = $2
			)
		ORDER BY nickname $1
		LIMIT $4::TEXT::INTEGER
	`
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
	f := models.Forum{}

	err := DB.pool.QueryRow(
			getForumSQL, 
			slug,
		).Scan(
			&f.Slug,
			&f.Title,
			&f.User,
			&f.Posts,
			&f.Threads,
	)

	if err != nil {
		return nil, ForumNotFound
	}

	return &f, nil
}

func CreateForumThreadDB(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		thread, err := GetThreadDB(t.Slug)
		if err != nil {
			return thread, ThreadIsExist
		}
	}

	err := DB.pool.QueryRow(
		createForumThreadSQL, 
		&t.Author,
		&t.Created, 
		&t.Forum,
		&t.Message,
		&t.Slug,
		&t.Title,
	).Scan(
		&t.Author,
		&t.Created, 
		&t.Forum,
		&t.ID,
		&t.Message,
		&t.Title,
	)

	if err != nil {
		switch err.(pgx.PgError).Code {
		case pgxErrNotNull:
			return nil, ForumOrAuthorNotFound
		case pgxErrForeignKey:
			return nil, ForumOrAuthorNotFound
		default:
			return nil, err
		}
	}

	return t, nil
}

func GetForumThreadsDB(slug string, limit, since, desc []byte) (*models.Threads, error) {
	var rows *pgx.Rows
	var err error

	if since != nil {
		dir := ">="
		ord := ""
		if bytes.Equal([]byte("true"), desc) {
			dir = "<="
			ord = "DESC"
		}
		rows, err = DB.pool.Query(
			getForumThreadsSQL,
			dir, 
			ord,
			slug, 
			since, 
			limit,
		)
	} else {
		ord := ""
		if bytes.Equal([]byte("true"), desc) {
			ord = "DESC"
		}
		rows, err = DB.pool.Query(
			getForumThreadsSinceSQL,
			ord, 
			slug, 
			limit,
		)
	}
	defer rows.Close()

	if err != nil {
		return nil, ForumNotFound
	}
	
	threads := models.Threads{}
	for rows.Next() {
		t := models.Thread{}
		err = rows.Scan(
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.ID,
			&t.Message,
			&t.Slug,
			&t.Title,
			&t.Votes,
		)
		threads = append(threads, &t)
	}

	if len(threads) == 0 {
		_, err := GetForumDB(slug)
		if err != nil {
			return nil, ForumNotFound
		}
	}	
	return &threads, nil
}

func GetForumUsersDB(slug string, limit, since, desc []byte) (*models.Users, error) {
	_, err := GetForumDB(slug)
	if err != nil {
		return nil, err
	}

	var rows *pgx.Rows
	//var err error

	if since != nil {
		dir := ">"
		ord := ""
		if bytes.Equal([]byte("true"), desc) {
			dir = "<"
			ord = "DESC"
		}
		rows, err = DB.pool.Query(
			getForumUsersSienceSQl,
			dir, 
			ord,
			slug, 
			since, 
			limit,
		)
	} else {
		ord := ""
		if bytes.Equal([]byte("true"), desc) {
			ord = "DESC"
		}
		rows, err = DB.pool.Query(
			getForumUsersSQl,
			ord, 
			slug, 
			limit,
		)
	}
	defer rows.Close()

	if err != nil {
		return nil, UserNotFound
	}
	
	users := models.Users{}
	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Email,
		)
		users = append(users, &u)
	}

	if len(users) == 0 {
		_, err := GetForumDB(slug)
		if err != nil {
			return nil, UserNotFound
		}
	}	
	return &users, nil
}
