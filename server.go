package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func Server(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error(w, http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("sessionID")
	data := map[string]interface{}{
		"IP": IP,
	}

	// log the login / register page
	if err != nil {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			fmt.Print("Error in finding the template")
			fmt.Println(err)
			Error(w, http.StatusInternalServerError)
			return
		}
		if tmpl.Execute(w, data) != nil {
			fmt.Print("Error in executing the template")
			fmt.Println(tmpl.Execute(w, data))
			Error(w, http.StatusInternalServerError)
			return
		}
		return
	}

	// log the homepage
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		Error(w, http.StatusInternalServerError)
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
		Error(w, http.StatusInternalServerError)
		return
	}

	data["ConnectedUser"] = connectedUser

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		fmt.Print("Error in finding the template")
		fmt.Println(err)
		Error(w, http.StatusInternalServerError)
		return
	}
	if tmpl.Execute(w, data) != nil {
		fmt.Print("Error in executing the template")
		fmt.Println(tmpl.Execute(w, data))
		Error(w, http.StatusInternalServerError)
		return
	}
}
