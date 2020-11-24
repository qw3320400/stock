package mysql

import (
	"database/sql"
	"stock-go/utils"

	_ "github.com/go-sql-driver/mysql"
)

const (
	connectionString = "root:123456@tcp(localhost:3306)/stock"
)

var (
	db *sql.DB
)

func Connect() error {
	var (
		err error
	)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return utils.Errorf(err, "sql.Open fail", err)
	}
	return nil
}

func Close() error {
	if db != nil {
		err := db.Close()
		if err != nil {
			return utils.Errorf(err, "db.Close fail", err)
		}
	}
	return nil
}
