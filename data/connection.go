package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(host, port, user, password, dbName string) (Db *sql.DB, err error) {
	Db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = Db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !databaseExists(Db, dbName) {
		if err := createDatabase(Db, dbName); err != nil {
			return nil, err
		}
	}

	if _, err := Db.Exec("USE " + dbName); err != nil {
		return nil, err
	}

	if !tableExists(Db, dbName, "your_table_name") {
		if err := createTable(Db); err != nil {
			return nil, err
		}
	}

	fmt.Println("Database and table are ready.")
	return Db, nil
}

func databaseExists(db *sql.DB, dbName string) bool {
	var exists string
	query := "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	err := db.QueryRow(query, dbName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return exists == dbName
}

func createDatabase(db *sql.DB, dbName string) error {
	_, err := db.Exec("CREATE DATABASE " + dbName)
	return err
}

func tableExists(db *sql.DB, dbName, tableName string) bool {
	var exists bool
	query := "SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ? LIMIT 1"
	err := db.QueryRow(query, dbName, tableName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	return exists
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE your_table_name (
			id INT AUTO_INCREMENT PRIMARY KEY,
			column_name_1 datatype_1,
			column_name_2 datatype_2,
			...
		)
	`
	_, err := db.Exec(query)
	return err
}
