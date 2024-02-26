package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func GiveCookie(w http.ResponseWriter, nickname string) {
	// Clear the sessionID cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	sessionID := uuid.New().String()
	expirationDate := time.Now().Add(24 * time.Hour)

	// Store sessionID in the database with user info
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE USERS SET cookie = ?, expiration = ? WHERE nickname = ?;",
		sessionID,
		expirationDate,
		nickname)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Set a cookie with the sessionID
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   sessionID,
		Expires: expirationDate,
	})
}
