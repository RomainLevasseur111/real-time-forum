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
	posts, err := getAllPosts()
	if err != nil {
		fmt.Println(err)
	}
	var connectedUser *USER

	if cookie != nil {
		connectedUser, _ = checkCookie(cookie)
	}

	// log the homepage
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		Error(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	data["ConnectedUser"] = connectedUser
	data["Posts"] = posts
	data["Categories"] = GetCategories()

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
