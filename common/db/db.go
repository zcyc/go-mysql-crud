package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:toor@tcp(localhost:3306)/dev")
	if err != nil {
		log.Fatal(err)
	}
}
