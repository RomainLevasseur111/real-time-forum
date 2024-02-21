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
		email TEXT NOT NULL UNIQUE,
		username VARCHAR(20),
		password TEXT,
		role VARCHAR(5),
		pfp TEXT,
		creationdate DATE,
		posts INTEGER DEFAULT 0,
		comments INTEGER DEFAULT 0,
		likes INTEGER DEFAULT 0,
		dislikes INTEGER DEFAULT 0,
		cookie TEXT,
		expiration DATE,
		session INTEGER DEFAULT 0,
		ask TEXT DEFAULT "false"
	);
	CREATE TABLE IF NOT EXISTS POSTS(
		userid INTEGER NOT NULL,
		postid INTEGER PRIMARY KEY AUTOINCREMENT,
		img TEXT DEFAULT NULL,
		commentid INTEGER DEFAULT NULL,
		username TEXT NOT NULL,
		role TEXT NOT NULL,
		category TEXT DEFAULT NULL,
		categoryB TEXT DEFAULT NULL,
		userpfp TEXT NOT NULL,
		content TEXT NOT NULL,
		postdate DATE NOT NULL,
		reasonreport TEXT DEFAULT "",
		reportid INTEGER,
		FOREIGN KEY(category) REFERENCES CATEGORIES(name),
		FOREIGN KEY(userid) REFERENCES USERS(id),
		FOREIGN KEY(commentid) REFERENCES POSTS(id),
		FOREIGN KEY(reportid) REFERENCES USERS(id)
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
	CREATE TABLE IF NOT EXISTS MESSAGES (
		messageid INTEGER PRIMARY KEY AUTOINCREMENT,
		receiverid INTERGER NOT NULL,
		sendername TEXT,
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
