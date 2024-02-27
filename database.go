package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// CreateDataBase function that creates the sqlite tables : USERS then POSTS
func CreateDatabase() {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	r := `
	CREATE TABLE IF NOT EXISTS USERS (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname VARCHAR(20) UNIQUE,
		age INTEGER NOT NULL,
		gender TEXT,
		firstname VARCHAR(20),
		lastname VARCHAR(20),
		email TEXT NOT NULL UNIQUE,
		password VARCHAR(30),
		pfp TEXT,
		creationdate DATE,
		cookie TEXT,
		expiration DATE
	);
	CREATE TABLE IF NOT EXISTS MESSAGES (
		messageid INTEGER PRIMARY KEY AUTOINCREMENT,
		sendername TEXT,
		receivername TEXT,
		content TEXT NOT NULL,
		FOREIGN KEY(receiverid) REFERENCES USERS(id)
	);
	 
	`
	_, err = db.Exec(r)
	if err != nil {
		log.Println("CREATE ERROR")
		fmt.Println(err)
	}
}
