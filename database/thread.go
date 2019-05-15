package database

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	"../models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
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

	// getThreadPosts
	getPostsSienceDescLimitTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 AND (path < (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
		ORDER BY path DESC
		LIMIT $3::TEXT::INTEGER
	`
	
	getPostsSienceDescLimitParentTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE path[1] IN (
			SELECT id
			FROM posts
			WHERE thread = $1 AND parent = 0 AND id < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
			ORDER BY id DESC
			LIMIT $3::TEXT::INTEGER
		)
		ORDER BY path	
	`
	getPostsSienceDescLimitFlatSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 AND id < $2::TEXT::INTEGER
		ORDER BY id DESC
		LIMIT $3::TEXT::INTEGER
	`
	getPostsSienceLimitTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
		ORDER BY path
		LIMIT $3::TEXT::INTEGER
	`
	getPostsSienceLimitParentTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE path[1] IN (
			SELECT id
			FROM posts
			WHERE thread = $1 AND parent = 0 AND id > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
			ORDER BY id LIMIT $3::TEXT::INTEGER
		)
		ORDER BY path
	`
	getPostsSienceLimitFlatSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 AND id > $2::TEXT::INTEGER
		ORDER BY id
		LIMIT $3::TEXT::INTEGER
	`
	// without sience
	getPostsDescLimitTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 
		ORDER BY path DESC
		LIMIT $2::TEXT::INTEGER
	`
	getPostsDescLimitParentTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 AND path[1] IN (
			SELECT path[1]
			FROM posts
			WHERE thread = $1
			GROUP BY path[1]
			ORDER BY path[1] DESC
			LIMIT $2::TEXT::INTEGER
		)
		ORDER BY path[1] DESC, path
	`
	getPostsDescLimitFlatSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1
		ORDER BY id DESC
		LIMIT $2::TEXT::INTEGER
	`
	getPostsLimitTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 
		ORDER BY path
		LIMIT $2::TEXT::INTEGER
	`
	getPostsLimitParentTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 AND path[1] IN (
			SELECT path[1] 
			FROM posts 
			WHERE thread = $1 
			GROUP BY path[1]
			ORDER BY path[1]
			LIMIT $2::TEXT::INTEGER
		)
		ORDER BY path
	`
	getPostsLimitFlatSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts
		WHERE thread = $1 
		ORDER BY id
		LIMIT $2::TEXT::INTEGER
	`
	// ThreadVote	
	getThreadVoteByIDSQL = `
		SELECT votes.voice, threads.id, threads.votes, u.nickname
		FROM (SELECT 1) s
		LEFT JOIN threads ON threads.id = $1
		LEFT JOIN "users" u ON u.nickname = $2
		LEFT JOIN votes ON threads.id = votes.thread AND u.nickname = votes.nickname
	`

	getThreadVoteBySlugSQL = `
		SELECT votes.voice, threads.id, threads.votes, u.nickname
		FROM (SELECT 1) s
		LEFT JOIN threads ON threads.slug = $1
		LEFT JOIN users as u ON u.nickname = $2
		LEFT JOIN votes ON threads.id = votes.thread AND u.nickname = votes.nickname
	`
	insertVoteSQL = `
		INSERT INTO votes (thread, nickname, voice) 
		VALUES ($1, $2, $3)
	`
	updateVoteSQL = `
		UPDATE votes SET 
		voice = $3
		WHERE thread = $1 
		AND nickname = $2
	`
	updateThreadVotesSQL = `
		UPDATE threads SET
		votes = $1
		WHERE id = $2
		RETURNING author, created, forum, "message" , slug, title, id, votes
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
	)
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
func parentExitsInOtherThread(parent int64, threadID int32) bool {
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

func parentNotExists(parent int64) bool {
	if parent == 0 {
		return false
	}

	var t int
	err := DB.pool.QueryRow(`SELECT id FROM posts WHERE id = $1`, parent,).Scan(&t); 
	
	if err != nil {
		return true
	}

	return false
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO posts (author, created, message, thread, parent, forum, path) VALUES ")
	queryBody := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (select currval(pg_get_serial_sequence('posts', 'id')))),"
	queryBodyEnd := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (select currval(pg_get_serial_sequence('posts', 'id'))))"
	for i, post := range *posts {
		if authorExists(post.Author) {
			return nil, UserNotFound
		}
		if parentExitsInOtherThread(post.Parent, thread.ID) || parentNotExists(post.Parent) {
			return nil, PostParentNotFound
		}

		// можно оптимизировать
		if i < len(*posts) - 1 {
			query.WriteString(fmt.Sprintf(queryBody, post.Author, created, post.Message, thread.ID, post.Parent, thread.Forum, post.Parent))
		} else {
			query.WriteString(fmt.Sprintf(queryBodyEnd, post.Author, created, post.Message, thread.ID, post.Parent, thread.Forum, post.Parent))
		}

	}
	query.WriteString("RETURNING author, created, forum, id, message, parent, thread")

	rows, err := DB.pool.Query(query.String()) 
	if err != nil {
		return nil, err
	}
	insertPosts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		_ = rows.Scan(
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.ID,
			&post.Message,
			&post.Parent,
			&post.Thread,
		)
		insertPosts = append(insertPosts, &post) 
	}
	return &insertPosts, nil
}

// НЕ ТЕСТИРОВАЛ
// /thread/{slug_or_id}/posts Сообщения данной ветви обсуждения
func GetThreadPostsDB(param, limit, since, sort, desc string) (*models.Posts, error) {
	thread, err := GetThreadDB(param)
	if err != nil {
		return nil, err
	}
	var rows *pgx.Rows

	if since != "" {
		if desc == "true" {
			switch sort {
			case "tree":
				rows, err = DB.pool.Query(getPostsSienceDescLimitTreeSQL, thread.ID, since, limit)
			case "parent_tree":
				rows, err = DB.pool.Query(getPostsSienceDescLimitParentTreeSQL, thread.ID, since, limit)
			default:
				rows, err = DB.pool.Query(getPostsSienceDescLimitFlatSQL, thread.ID, since, limit)
			}
		} else {
			switch sort {
			case "tree":
				rows, err = DB.pool.Query(getPostsSienceLimitTreeSQL, thread.ID, since, limit)
			case "parent_tree":
				rows, err = DB.pool.Query(getPostsSienceLimitParentTreeSQL, thread.ID, since, limit)
			default:
				rows, err = DB.pool.Query(getPostsSienceLimitFlatSQL, thread.ID, since, limit)
			}
		}
	} else {
		if desc == "true" {
			switch sort {
			case "tree":
				rows, err = DB.pool.Query(getPostsDescLimitTreeSQL, thread.ID, limit)
			case "parent_tree":
				rows, err = DB.pool.Query(getPostsDescLimitParentTreeSQL, thread.ID, limit)
			default:
				rows, err = DB.pool.Query(getPostsDescLimitFlatSQL, thread.ID, limit)
			}
		} else {
			switch sort {
			case "tree":
				rows, err = DB.pool.Query(getPostsLimitTreeSQL, thread.ID, limit)
			case "parent_tree":
				rows, err = DB.pool.Query(getPostsLimitParentTreeSQL, thread.ID, limit)
			default:
				rows, err = DB.pool.Query(getPostsLimitFlatSQL, thread.ID, limit)
			}
		}
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}

		if err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return &posts, nil
}

// УБРАТЬ КОСТЫЛИ
// НЕ ТЕСТИРОВАЛ
// /thread/{slug_or_id}/vote Проголосовать за ветвь обсуждения
func MakeThreadVoteDB(vote *models.Vote, param string) *models.Thread {
	var err error
	prevVoice := &pgtype.Int4{}
	threadID := &pgtype.Int4{}
	threadVotes := &pgtype.Int4{}
	userNickname := &pgtype.Varchar{}

	if isNumber(param) {
		id, _ := strconv.Atoi(param)
		err = DB.pool.QueryRow(getThreadVoteByIDSQL, id, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	} else {
		err = DB.pool.QueryRow(getThreadVoteBySlugSQL, param, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	}
	if err != nil {
		return nil
	}
	if threadID.Status != pgtype.Present || userNickname.Status != pgtype.Present {
		return nil
	}
	var prevVoiceInt int32
	if prevVoice.Status == pgtype.Present {
		prevVoiceInt = int32(prevVoice.Int)
		_, err = DB.pool.Exec(updateVoteSQL, threadID.Int, userNickname.String, vote.Voice)
	} else {
		_, err = DB.pool.Exec(insertVoteSQL, threadID.Int, userNickname.String, vote.Voice)
	}
	newVotes := threadVotes.Int + (int32(vote.Voice) - prevVoiceInt)
	if err != nil {
		return nil
	}
	thread := &models.Thread{}
	slugNullable := &pgtype.Varchar{}
	err = DB.pool.QueryRow(
		updateThreadVotesSQL,
		newVotes,
		threadID.Int,
	).Scan(
		&thread.Author,
		&thread.Created,
		thread.Forum,
		&thread.Message,
		slugNullable,
		&thread.Title,
		&thread.ID,
		&thread.Votes,
	)
	thread.Slug = slugNullable.String
	if err != nil {
		return nil
	}

	return thread
}