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
	
	// getPostsSienceDescLimitParentTreeSQL = `
	// 	SELECT id, author, parent, message, forum, thread, created
	// 	FROM posts
	// 	WHERE path[1] IN (
	// 		SELECT id
	// 		FROM posts
	// 		WHERE thread = $1 AND parent = 0 AND id < (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
	// 		ORDER BY id DESC
	// 		LIMIT $3::TEXT::INTEGER
	// 	)
	// 	ORDER BY path
	// `

	getPostsSienceDescLimitParentTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts p
		WHERE p.thread = $1 and p.path[1] IN (
			SELECT p2.path[1]
			FROM posts p2
			WHERE p2.thread = $1 AND p2.parent = 0 and p2.path[1] < (SELECT p3.path[1] from posts p3 where p3.id = $2)
			ORDER BY p2.path DESC
			LIMIT $3
		)
		ORDER BY p.path[1] DESC, p.path[2:]
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
	// getPostsSienceLimitParentTreeSQL = `
	// 	SELECT id, author, parent, message, forum, thread, created
	// 	FROM posts
	// 	WHERE path[1] IN (
	// 		SELECT id
	// 		FROM posts
	// 		WHERE thread = $1 AND parent = 0 AND id > (SELECT path[1] FROM posts WHERE id = $2::TEXT::INTEGER)
	// 		ORDER BY id 
	// 		LIMIT $3::TEXT::INTEGER
	// 	)
	// 	ORDER BY path
	// `
	getPostsSienceLimitParentTreeSQL = `
		SELECT id, author, parent, message, forum, thread, created
		FROM posts p
		WHERE p.thread = $1 and p.path[1] IN (
			SELECT p2.path[1]
			FROM posts p2
			WHERE p2.thread = $1 AND p2.parent = 0 and p2.path[1] > (SELECT p3.path[1] from posts p3 where p3.id = $2::TEXT::INTEGER)
			ORDER BY p2.path
			LIMIT $3::TEXT::INTEGER
		)
		ORDER BY p.path
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
		RETURNING author, created, forum, "message", slug, title, id, votes
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
		// возможно можно сократить
		id, _ := strconv.Atoi(param)
		err = DB.pool.QueryRow(
			getThreadIdSQL,
			id,
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
	} else {
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
	}

	if err != nil {
		return nil, ThreadNotFound
	}

	return &thread, nil
}

// /thread/{slug_or_id}/details Обновление ветки
func UpdateThreadDB(thread *models.ThreadUpdate, param string) (*models.Thread, error) {
	threadFound, err := GetThreadDB(param)
	if err != nil {
		return nil, PostNotFound
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
	var t int64
	rows := DB.pool.QueryRow(`
		SELECT id
		FROM posts
		WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)`,
		parent, 
		threadID,
	)

	err := rows.Scan(&t)	
	if err != nil && err.Error() == "no rows in result set" {
		return false
	}
	return true
}

func parentNotExists(parent int64) bool {
	if parent == 0 {
		return false
	}

	var t int64
	err := DB.pool.QueryRow(`SELECT id FROM posts WHERE id = $1`, parent).Scan(&t); 
	
	if err != nil {
		return true
	}

	return false
}

func checkPost(p *models.Post, t *models.Thread) error {
	if authorExists(p.Author) {
		return UserNotFound
	}
	if parentExitsInOtherThread(p.Parent, t.ID) || parentNotExists(p.Parent) {
		return PostParentNotFound
	}
	return nil
}

// thread/{slug_or_id}/create Создание новых постов
func CreateThreadDB(posts *models.Posts, param string) (*models.Posts, error) {
	thread, err := GetThreadDB(param)
	if err != nil {
		return nil, err
	}

	postsNumber := len(*posts)

	if postsNumber == 0 {
		return posts, nil
	}

	dateTimeTemplate := "2006-01-02 15:04:05"
	created := time.Now().Format(dateTimeTemplate)
	query := strings.Builder{}
	query.WriteString("INSERT INTO posts (author, created, message, thread, parent, forum, path) VALUES ")
	queryBody := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (SELECT last_value FROM posts_id_seq)),"
	for i, post := range *posts {
		err = checkPost(post, thread)
		if err != nil {
			return nil, err
		}

		temp := fmt.Sprintf(queryBody, post.Author, created, post.Message, thread.ID, post.Parent, thread.Forum, post.Parent)
		// удаление запятой в конце queryBody для последнего подзапроса
		if i == postsNumber - 1 {
			temp = temp[:len(temp) - 1]
		}
		query.WriteString(temp)
	}
	query.WriteString("RETURNING author, created, forum, id, message, parent, thread")

	rows, err := DB.pool.Query(query.String())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	insertPosts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		rows.Scan(
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
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}	
	return &insertPosts, nil
}

var queryPostsWithSience = map[string]map[string]string {
	"true": map[string]string {
		"tree": getPostsSienceDescLimitTreeSQL,
		"parent_tree": getPostsSienceDescLimitParentTreeSQL,
		"flat": getPostsSienceDescLimitFlatSQL,
	},
	"false": map[string]string {
		"tree": getPostsSienceLimitTreeSQL,
		"parent_tree": getPostsSienceLimitParentTreeSQL,
		"flat": getPostsSienceLimitFlatSQL,
	},
}

var queryPostsNoSience = map[string]map[string]string {
	"true": map[string]string {
		"tree": getPostsDescLimitTreeSQL,
		"parent_tree": getPostsDescLimitParentTreeSQL,
		"flat": getPostsDescLimitFlatSQL,
	},
	"false": map[string]string {
		"tree": getPostsLimitTreeSQL,
		"parent_tree": getPostsLimitParentTreeSQL,
		"flat": getPostsLimitFlatSQL,
	},
}

// /thread/{slug_or_id}/posts Сообщения данной ветви обсуждения
func GetThreadPostsDB(param, limit, since, sort, desc string) (*models.Posts, error) {
	// tmp, _ := strconv.Atoi(limit)
	// tmp += 1;
	// limit = strconv.Itoa(tmp)
	thread, err := GetThreadDB(param)
	if err != nil {
		return nil, ForumNotFound
	}

	var rows *pgx.Rows

	if since != "" {
		query := queryPostsWithSience[desc][sort]
		rows, err = DB.pool.Query(query, thread.ID, since, limit)
	} else {
		query := queryPostsNoSience[desc][sort]
		rows, err = DB.pool.Query(query, thread.ID, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	// fmt.Println(len(posts))
	return &posts, nil
}

// УБРАТЬ КОСТЫЛИ
// /thread/{slug_or_id}/vote Проголосовать за ветвь обсуждения
func MakeThreadVoteDB(vote *models.Vote, param string) (*models.Thread, error) {
	var err error

	var thrID int
	if isNumber(param) {
		id, _ := strconv.Atoi(param)
		err = DB.pool.QueryRow(`SELECT id FROM threads WHERE id = $1`, id).Scan(&thrID)
	} else {
		err = DB.pool.QueryRow(`SELECT id FROM threads WHERE slug = $1`, param).Scan(&thrID)
	}	
	if err != nil {
		return nil, ForumNotFound
	}
	var uNick string
	err = DB.pool.QueryRow(`SELECT nickname FROM users WHERE nickname = $1`, vote.Nickname).Scan(&uNick)
	if err != nil {
		return nil, UserNotFound
	}
	prevVoice := &pgtype.Int4{}
	threadID := &pgtype.Int4{}
	threadVotes := &pgtype.Int4{}
	userNickname := &pgtype.Varchar{}
	err = DB.pool.QueryRow(getThreadVoteByIDSQL, thrID, vote.Nickname).Scan(prevVoice, threadID, threadVotes, userNickname)
	if err != nil || threadID.Status != pgtype.Present || userNickname.Status != pgtype.Present {
		return nil, err
	}

	var newVotes int32
	if prevVoice.Status == pgtype.Present {
		_, err = DB.pool.Exec(updateVoteSQL, threadID.Int, userNickname.String, vote.Voice)
		newVotes = threadVotes.Int + (vote.Voice - prevVoice.Int)
	} else {
		_, err = DB.pool.Exec(insertVoteSQL, threadID.Int, userNickname.String, vote.Voice)
		newVotes = threadVotes.Int + vote.Voice
	}
	if err != nil {
		return nil, err
	}
	thread := &models.Thread{}
	err = DB.pool.QueryRow(
		updateThreadVotesSQL,
		newVotes,
		threadID.Int,
	).Scan(
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.ID,
		&thread.Votes,
	)
	if err != nil {
		return nil, err
	}

	return thread, nil
}