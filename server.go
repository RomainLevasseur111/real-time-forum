package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func Server(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		// log the register / login page
		return
	}
	// log the homepage

	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var connectedUser USER

	row := db.QueryRow(`SELECT id, nickname, age, gender, firstname, lastname, email, pfp, creationdate FROM USERS WHERE cookie = ?;`, cookie.Value)
	err = row.Scan(
		&connectedUser.Id,
		&connectedUser.NickName,
		&connectedUser.Age,
		&connectedUser.Gender,
		&connectedUser.FirstName,
		&connectedUser.LastName,
		&connectedUser.Email,
		&connectedUser.Pfp,
		&connectedUser.CreationDate,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := map[string]interface{}{
		"IP":            IP,
		"ConnectedUser": connectedUser,
	}

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		fmt.Print("Error in finding the template")
		fmt.Println(err)
		return
	}
	if tmpl.Execute(w, data) != nil {
		fmt.Print("Error in executing the template")
		fmt.Println(tmpl.Execute(w, data))
		return
	}
}
