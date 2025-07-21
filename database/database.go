package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	itemModel "test.com/events/model/itemModel"
)

var DB *sql.DB

func Connect() error {
	var dsn string = os.Getenv("DB_URL")

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	itemModel.StartDb(DB)
	fmt.Println("Successfully connected to database!")
	return nil
}
