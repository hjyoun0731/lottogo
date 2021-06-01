package api

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	UserId       int
	UserName     string
	UserPassword string
	Created      string
	Updated      string
}

func NewDb() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/lottoDb")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func CloseDb(db *sql.DB) {
	db.Close()
}
