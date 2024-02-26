package main

import (
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
	// log the homepage with the cookie
	fmt.Println(cookie)

	data := map[string]interface{}{
		"IP": IP,
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
