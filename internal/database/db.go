package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// db *sql.DB
	once     sync.Once
	instance *sql.DB
)

// Returns a single instance of sql.DB
func GetInstance() *sql.DB {
	log.Println("Getting Db instance...")
	once.Do(func() {
		instance = dbConnect()
		log.Println("Db Instance created...")
	})
	return instance
}

// Connect to database
func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "chatapp.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Close db connection
func dbClose() {
	defer instance.Close()
}
