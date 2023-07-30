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

	tableNames := []string{"users", "weather_history"}
	for _, tableName := range tableNames {
		if !tableExists(Db, dbName, tableName) {
			if err := createTable(Db, tableName); err != nil {
				return nil, err
			}
		}
	}

	fmt.Println("Connected to database")
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

func createTable(db *sql.DB, tableName string) error {

	var query string
	switch tableName {
	case "users":
		query = `
		CREATE TABLE users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			date_of_birth DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB;
		`
	case "weather_history":
		query = `
		CREATE TABLE weather_history (
		  id INT NOT NULL AUTO_INCREMENT,
		  city_name VARCHAR(255) NOT NULL,
		  user_id INT NOT NULL,
		  coord_lon FLOAT NOT NULL,
		  coord_lat FLOAT NOT NULL,
		  weather_id INT NOT NULL,
		  weather_main VARCHAR(255) NOT NULL,
		  weather_description VARCHAR(255) NOT NULL,
		  weather_icon VARCHAR(255) NOT NULL,
		  base VARCHAR(255) NOT NULL,
		  temp FLOAT NOT NULL,
		  feels_like FLOAT NOT NULL,
		  temp_min FLOAT NOT NULL,
		  temp_max FLOAT NOT NULL,
		  pressure INT NOT NULL,
		  humidity INT NOT NULL,
		  visibility INT NOT NULL,
		  wind_speed FLOAT NOT NULL,
		  wind_deg INT NOT NULL,
		  clouds_all INT NOT NULL,
		  dt INT NOT NULL,
		  sys_type INT NOT NULL,
		  sys_id INT NOT NULL,
		  sys_country VARCHAR(255) NOT NULL,
		  sys_sunrise INT NOT NULL,
		  sys_sunset INT NOT NULL,
		  timezone INT NOT NULL,
		  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  PRIMARY KEY (id),
		  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
		) ENGINE=InnoDB;`

	}

	_, err := db.Exec(query)
	return err
}
