package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Server struct {
	Database *Database
}
type Database struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Name       string
	Connection *sql.DB
}

func (db *Database) Connect() error {

	db.Host = os.Getenv("SIJITPOS_DB_HOST")
	if db.Host == "" {
		db.Host = "localhost"
	}
	db.Port = os.Getenv("SIJILPOS_DB_PORT")
	if db.Port == "" {
		db.Port = "5432"
	}

	db.Username = os.Getenv("SIJILPOS_DB_USER")
	if db.Username == "" {
		db.Username = "postgres"
	}
	db.Password = os.Getenv("SIJILPOS_DB_PASSWORD")
	if db.Password == "" {
		log.Fatalln("Database password is not set in environment variables")
	}
	db.Name = os.Getenv("SIJILPOS_DB_NAME")
	if db.Name == "" {
		db.Name = "sijilposdb"
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.Username, db.Password, db.Name)
	var err error
	db.Connection, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return db.Connection.Ping()
}
