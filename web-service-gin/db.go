package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"os"
)

func getDbProvider() (*sql.DB, error) {
	var db *sql.DB
	cfg := mysql.Config{
		User:   os.Getenv("DB-USER"),
		Passwd: os.Getenv("DB-PASS"),
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "world",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return db, err
	}
	pingErr := db.Ping()
	if pingErr != nil {
		return db, pingErr
	}
	return db, nil
}
