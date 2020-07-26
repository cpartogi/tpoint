package db

import (
	"database/sql"

	"github.com/cpartogi/tpoint/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

var err error

func Init() {
	conf := config.GetConfig()

	connString := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connString)
	if err != nil {
		panic("connection error")
	}
}

func CreateCon() *sql.DB {
	return db
}
