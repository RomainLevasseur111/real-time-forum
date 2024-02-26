package main

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Check if Email already Exist
func EmailAlreadyExist(email string) bool {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT email FROM USERS WHERE email = ?", email)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var emailExists string

	for rows.Next() {
		rows.Scan(&emailExists)
	}
	return emailExists == ""
}

// Check if nickname already exists
func NicknameAlreadyExists(nickname string) bool {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT nickname FROM USERS WHERE nickname = ?", nickname)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var nicknameExists string

	for rows.Next() {
		rows.Scan(&nicknameExists)
	}
	return nicknameExists == ""
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
