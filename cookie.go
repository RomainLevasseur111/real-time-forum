package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func checkCookie() {
	/*
			// Get the sessionID from the cookie
			sessionID, err := r.Cookie("sessionID")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			// Verify sessionID in the database

			// If valid, serve the dashboard content
			// Else, redirect to login
		}
	*/
}

//
// LOG-OUT

/*
	// Invalidate session by removing it from the database

	// Clear the sessionID cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	// Redirect to login page probably
*/
