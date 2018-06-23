package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qichengzx/shortme/config"
)

var (
	db      *sql.DB
	DSN     string
	err     error
	DB_ADDR string
	DB_NAME string
	DB_USER string
	DB_PASS string
)

func init() {
	DB_ADDR, err = config.GetByBlock("mysql", "mysql.addr")
	DB_NAME, err = config.GetByBlock("mysql", "mysql.db")
	DB_USER, err = config.GetByBlock("mysql", "mysql.user")
	DB_PASS, err = config.GetByBlock("mysql", "mysql.password")

	DSN = DB_USER + ":" + DB_PASS + "@" + DB_ADDR + "/" + DB_NAME + "?charset=utf8&loc=Asia%2FShanghai&parseTime=true"

	db, err = sql.Open("mysql", DSN)
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
}
