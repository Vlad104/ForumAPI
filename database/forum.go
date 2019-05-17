package database

import (
	"fmt"
	"../models"
	"github.com/jackc/pgx"
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
		INTO threads (author, created, message, title, slug, forum)
		VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6)) 
		RETURNING author, created, forum, id, message, title
	`
	getForumThreadsSinceSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
		ORDER BY created
		LIMIT $3::TEXT::INTEGER
	`
	getForumThreadsDescSinceSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
		ORDER BY created DESC
		LIMIT $3::TEXT::INTEGER
	`
	getForumThreadsSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1
		ORDER BY created
		LIMIT $2::TEXT::INTEGER
	`
	getForumThreadsDescSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1
		ORDER BY created DESC
		LIMIT $2::TEXT::INTEGER
	`
	getForumUsersSienceSQl = `
	SELECT nickname, fullname, about, email
	FROM users
	WHERE nickname IN (
			SELECT author FROM threads WHERE forum = $1
			UNION
			SELECT author FROM posts WHERE forum = $1
		) 
		AND LOWER(nickname) > LOWER($1::TEXT)
	ORDER BY nickname
	LIMIT $3::TEXT::INTEGER
	`
	getForumUsersDescSienceSQl = `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT author FROM threads WHERE forum = $1
				UNION
				SELECT author FROM posts WHERE forum = $1
			) 
			AND LOWER(nickname) < LOWER($2::TEXT)
		ORDER BY nickname DESC
		LIMIT $3::TEXT::INTEGER
	`
	getForumUsersSQl = `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT author FROM threads WHERE forum = $1
				UNION
				SELECT author FROM posts WHERE forum = $1
			)
		ORDER BY nickname
		LIMIT $2::TEXT::INTEGER
	`
	getForumUsersDescSQl = `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT author FROM threads WHERE forum = $1
				UNION
				SELECT author FROM posts WHERE forum = $1
			)
		ORDER BY nickname DESC
		LIMIT $2::TEXT::INTEGER
	`
)

// /forum/create Создание форума
func CreateForumDB(f *models.Forum) (*models.Forum, error)  {
	err := DB.pool.QueryRow(
		createForumSQL,
		&f.Slug,
		&f.Title,
		&f.User,
	).Scan(&f.User)

	switch ErrorCode(err) {
	case pgxOK:
		return f, nil
	case pgxErrUnique:
		forum, _ := GetForumDB(f.Slug)
		return forum, ForumIsExist
	case pgxErrNotNull:
		return nil, UserNotFound
	default:
		return nil, err
	}
}

// /forum/{slug}/details Получение информации о форуме
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

// /forum/{slug}/create Создание ветки
func CreateForumThreadDB(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		thread, err := GetThreadDB(t.Slug)
		if err == nil {
			return thread, ThreadIsExist
		}
	}

	err := DB.pool.QueryRow(
		createForumThreadSQL, 
		&t.Author,
		&t.Created, 
		&t.Message,
		&t.Title,
		&t.Slug,
		&t.Forum,
	).Scan(
		&t.Author,
		&t.Created, 
		&t.Forum,
		&t.ID,
		&t.Message,
		&t.Title,
	)
	
	switch ErrorCode(err) {
	case pgxOK:
		return t, nil
	case pgxErrNotNull:
		return nil, ForumOrAuthorNotFound //UserNotFound
	case pgxErrForeignKey:
		return nil, ForumOrAuthorNotFound //ForumIsExist
	default:
		return nil, err
	}
}

// /forum/{slug}/threads Список ветвей обсужления форума
func GetForumThreadsDB(slug, limit, since, desc string) (*models.Threads, error) {
	var rows *pgx.Rows
	var err error

	if since != "" {
		if desc == "true" {
			rows, err = DB.pool.Query(
				getForumThreadsDescSinceSQL,
				slug,
				since,
				limit,
			)
		} else {
			rows, err = DB.pool.Query(
				getForumThreadsSinceSQL,
				slug,
				since,
				limit,
			)
		}
	} else {
		if desc == "true" {
			rows, err = DB.pool.Query(
				getForumThreadsDescSQL,
				slug,
				limit,
			)
		} else {
			rows, err = DB.pool.Query(
				getForumThreadsSQL,
				slug,
				limit,
			)
		}
	}
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
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

// /forum/{slug}/users Пользователи данного форума
func GetForumUsersDB(slug string, limit, since, desc string) (*models.Users, error) {
	var rows *pgx.Rows
	var err error

	if since != "" {
		if desc == "true" {
			rows, err = DB.pool.Query(
				getForumUsersDescSienceSQl,
				slug,
				since,
				limit,
			)
		} else {
			rows, err = DB.pool.Query(
				getForumUsersSienceSQl,
				slug,
				since,
				limit,
			)
		}
	} else {
		if desc == "true" {
			rows, err = DB.pool.Query(
				getForumUsersDescSQl,
				slug,
				limit,
			)
		} else {
			rows, err = DB.pool.Query(
				getForumUsersSQl,
				slug,
				limit,
			)
		}
	}
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return nil, ForumNotFound
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
