package mysql

import (
	"database/sql"
	"fmt"
	"stock-go/utils"

	_ "github.com/go-sql-driver/mysql"
)

const (
	connectionString = "root:123456@tcp(localhost:3306)/stock"
)

var (
	db *sql.DB
)

func GetConnection() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	tmpDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, utils.Errorf(err, "sql.Open fail", err)
	}
	db = tmpDB
	return db, nil
}

func CloseConnection() error {
	if db == nil {
		return nil
	}
	tmpDB := db
	db = nil
	err := tmpDB.Close()
	if err != nil {
		utils.LogErr(fmt.Sprintf("db.Close fail %s", err))
	}
	return nil
}
