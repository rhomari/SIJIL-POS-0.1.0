package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Server struct {
	Database *Database
}
type Database struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Name       string
	Connection *sql.DB
}

func (db *Database) Connect() error {

	db.Host = "localhost"
	db.Port = 5432
	db.Username = "user"
	db.Password = "password"
	db.Name = "dbname"
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.Username, db.Password, db.Name)
	var err error
	db.Connection, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return db.Connection.Ping()
}
