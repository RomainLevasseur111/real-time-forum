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
		nickname VARCHAR(9) UNIQUE,
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
		sendername VARCHAR(20),
		receivername VARCHAR(20),
		date TEXT,
		pfp TEXT,
		content TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS POSTS(
		userid INTEGER NOT NULL,
		postid INTEGER PRIMARY KEY AUTOINCREMENT,
		commentid INTEGER DEFAULT NULL,
		username TEXT NOT NULL,
		category TEXT DEFAULT NULL,
		categoryB TEXT DEFAULT NULL,
		userpfp TEXT NOT NULL,
		content TEXT NOT NULL,
		postdate DATE NOT NULL,
		FOREIGN KEY(category) REFERENCES CATEGORIES(name),
		FOREIGN KEY(userid) REFERENCES USERS(id),
		FOREIGN KEY(commentid) REFERENCES POSTS(id)
	);
	CREATE TABLE IF NOT EXISTS CATEGORIES(
		name TEXT UNIQUE,
		posts INTEGER
	);
	CREATE TABLE IF NOT EXISTS LikesDislikes (
		id INTEGER PRIMARY KEY,
		userid INTEGER NOT NULL,
		postid INTEGER NOT NULL,
		likestatus TEXT NOT NULL CHECK(likestatus IN ('like', 'dislike')),
		FOREIGN KEY(userid) REFERENCES USERS(id),
		FOREIGN KEY(postid) REFERENCES POSTS(id)
	);

	`
	_, err = db.Exec(r)
	if err != nil {
		log.Println("CREATE ERROR")
		fmt.Println(err)
	}
}
