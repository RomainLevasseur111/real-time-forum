package main

import (
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Check login datas in the database

	// If authentication is successful:
	GiveCookie(w, "nickname here")

	// Redirect to another page
}
