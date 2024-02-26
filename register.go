package main

import (
	"database/sql"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// error 405 method not allowed
	}

	// Check if registeration datas are correct
	nickname := r.FormValue("nickname")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	email := r.FormValue("email")
	password := r.FormValue("password")

	var errMsg string
	errMsg = CheckName(nickname)
	if errMsg != "" {
		return
	}
	if !NicknameAlreadyExists(nickname) {
		errMsg = "Nickname already exists"
		return
	}

	errMsg = CheckName(firstname)
	if errMsg != "" {
		return
	}

	errMsg = CheckName(lastname)
	if errMsg != "" {
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		errMsg = "Invalid age format"
		return
	}
	if ageInt < 0 {
		errMsg = "Age can't be negative"
		return
	}

	if gender != "man" && gender != "woman" && gender != "other" {
		errMsg = "Invalid gender format"
		return
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		errMsg = "Invalid email address"
		return
	}
	if !EmailAlreadyExist(email) {
		errMsg = "Email already exists"
		return
	}

	if len(password) > 30 {
		errMsg = "Password can't have more than 30 characters"
		return
	}
	if len(password) < 8 {
		errMsg = "Password can't have less than 8 characters"
		return
	}
	psw, err := HashPassword(password)
	if err != nil {
		errMsg = "Error hashing your password, try another password"
		return
	}

	// If regisgtration is successful:
	// Store datas in database

	// Create a session cookie
	Successful_Login(w)

	// Redirect to another page
}

func CheckName(name string) string {
	if name == "" {
		return "Name can't be null"
	}
	if len(name) > 20 {
		return "Name can't have more than 20 characters"
	}
	if len(name) < 4 {
		return "Name can't have less than 4 characters"
	}
	if strings.Contains(name, " ") {
		return "Name can't have a space in it"
	}

	return ""
}

// Hash a given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check if Email already Exist
func EmailAlreadyExist(email string) bool {
	db, _ := sql.Open(DRIVER, DB)
	defer db.Close()
	rows, _ := db.Query("SELECT email FROM USERS WHERE email = ?", email)

	var emailExists string

	for rows.Next() {
		rows.Scan(&emailExists)
	}
	return emailExists == ""
}

// Check if nickname already exists
func NicknameAlreadyExists(nickname string) bool {
	db, _ := sql.Open(DRIVER, DB)
	defer db.Close()
	rows, _ := db.Query("SELECT username FROM USERS WHERE nickname = ?", nickname)

	var nicknameExists string

	for rows.Next() {
		rows.Scan(&nicknameExists)
	}
	return nicknameExists == ""
}
