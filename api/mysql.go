package api

import (
	"database/sql"
	"log"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Created  int    `json:"created"`
	Updated  int    `json:"updated"`
}

func NewTable() error {
	db := NewDb()
	defer db.Close()
	_,err := db.Exec(`create table if not exists user_table ( \
				id INT AUTO_INCREMENT, \
				name CHAR(32) DEFAULT NULL, \
				password CHAR(32) DEFAULT NULL, \
				created TIMESTAMP, \
				updated TIMESTAMP, \
				PRIMARY KEY(id) \
			)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)
	if err != nil {
		log.Println("create table fail.")
		return errors.New("create table fail")
	}
	return nil
}

func NewDb() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(mysql)/lottodb")
	if err != nil {
		log.Println(err)
	}

	return db
}

func CloseDb(db *sql.DB) {
	db.Close()
}

func queryPw(db *sql.DB, name string) (int, string) {
	var pw string
	var id int
	err := db.QueryRow("SELECT id, password FROM user_table where name = ?", name).Scan(&id, &pw)
	if err != nil {
		log.Println(err)
	}
	return id, pw
}
