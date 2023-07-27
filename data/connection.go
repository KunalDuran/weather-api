package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(host, port, user, password string) (Db *sql.DB, err error) {
	Db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/weather")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = Db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return Db, nil
}
