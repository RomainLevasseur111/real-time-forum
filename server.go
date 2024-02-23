package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func Server(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
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
