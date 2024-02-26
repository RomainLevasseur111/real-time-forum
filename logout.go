package main

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Error(w, http.StatusMethodNotAllowed)
		return
	}

	// Clear the sessionID cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	// Redirect to homepage
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
