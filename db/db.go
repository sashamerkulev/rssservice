package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DBOpen() error {
	mysql, err := sql.Open("mysql", "news:News,News@/dbnews?parseTime=true")
	if err != nil {
		return err
	}
	DB = mysql
	return nil
}

func DBClose() error {
	err := DB.Close()
	if err != nil {
		return err
	}
	return nil
}
