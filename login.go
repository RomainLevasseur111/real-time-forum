package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Check login datas in the database
	if r.Method != "POST" {
		Error(w, http.StatusMethodNotAllowed)
		return
	}
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		Error(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var pass string
	name := r.FormValue("login")
	pwd := r.FormValue("password")
	if !EmailAlreadyExist(name) {
		rows, err := db.Query("SELECT nickname FROM USERS WHERE email = ?", name)
		if err != nil {
			fmt.Println(err)
			Error(w, http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			rows.Scan(&name)
		}
	}
	if !NicknameAlreadyExists(name) {
		rows, err := db.Query("SELECT password FROM USERS WHERE nickname = ?", name)
		if err != nil {
			fmt.Println(err)
			Error(w, http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			rows.Scan(&pass)
		}
	}
	if CheckPasswordHash(pwd, pass) {
		GiveCookie(w, name)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
