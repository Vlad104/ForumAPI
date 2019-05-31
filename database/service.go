package database

import (
	"../models"
)

const (
	getStatusSQL = `
		SELECT *
		FROM (SELECT COUNT(*) FROM "users") as "users"
		CROSS JOIN (SELECT COUNT(*) FROM "threads") as threads
		CROSS JOIN (SELECT COUNT(*) FROM "forums") as forums
		CROSS JOIN (SELECT COUNT(*) FROM "posts") as posts
	`
	clearSQL = `
		TRUNCATE users, forums, threads, posts, votes;
	`
)

// /handlers/status Получение инфомарции о базе данных
func GetStatusDB() *models.Status {
	status := &models.Status{}
	DB.pool.QueryRow(
		getStatusSQL,
	).Scan(
		&status.User,
		&status.Thread,
		&status.Forum,
		&status.Post,
	)
	return status
}

// /handlers/clear Очистка всех данных в базе
func ClearDB() {
	DB.pool.Exec(clearSQL)
}