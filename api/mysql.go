package api

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Created  int    `json:"created"`
	Updated  int    `json:"updated"`
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

func queryPw(db *sql.DB, name string) (int, string) {
	var pw string
	var id int
	err := db.QueryRow("SELECT id, name FROM User_Table where name = ?", name).Scan(&id, &pw)
	if err != nil {
		log.Fatal(err)
	}
	return id, pw
}
