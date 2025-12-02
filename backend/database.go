package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "./santa.db")
	if err != nil {
		log.Fatal()
	}

	// Tworzenie tabeli participants
	createTableSQL := `CREATE TABLE IF NOT EXISTS participants (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"name" TEXT,
	"email" TEXT UNIQUE,
	"preferences" TEXT
	);`

	log.Println("Inicjalizacja bazy danych...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	log.Println("Baza danych gotowa (santa.db)")
}
