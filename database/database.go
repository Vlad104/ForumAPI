package database

import (
	//"database/sql"
	"github.com/jackc/pgx"
)

type DataBase struct {
	pool *pgx.ConnPool
}

var DB DataBase

/*
func NewDataBase() *DataBase, error {
	var db DataBase
	db.Connect()
	return &db, err
}
*/

func (db *DataBase) Connect() error {
	runtimeParams := make(map[string] string)
	runtimeParams["application_name"] = "dz"
	conConfig := pgx.ConnConfig {
		// Host: 			"localhost",
		Host: 			"127.0.0.1",
		Port: 			5432,
		Database: 		"mydb",
		User: 			"postgres",
		Password: 		"postgres",
		TLSConfig: 		nil,
		UseFallbackTLS: false,
		RuntimeParams: 	runtimeParams,
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     conConfig,
		MaxConnections: 50,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	p, err := pgx.NewConnPool(poolConfig)
	db.pool = p
	
	return err
}

func ErrorCode(err error) (string) {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return pgxOK
	}
	return pgerr.Code
}
