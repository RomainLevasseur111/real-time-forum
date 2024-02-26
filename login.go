<<<<<<< HEAD
package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Check login datas in the database

	// If authentication is successful:
	Successful_Login(w)

	// Redirect to another page
}

func Successful_Login(w http.ResponseWriter) {
	sessionID := uuid.New().String()
	expirationDate := time.Now().Add(24 * time.Hour)

	// Store sessionID in the database with user info

	// Set a cookie with the sessionID
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   sessionID,
		Expires: expirationDate,
	})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
=======
package main
>>>>>>> romain
