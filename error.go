package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// log error page
func Error(w http.ResponseWriter, statusCode int) {
	var msg string
	var name string

	switch statusCode {

	case http.StatusNotFound:
		name = "Not Found"
		msg = "This page does not exist"
	case http.StatusMethodNotAllowed:
		name = "Method not allowed"
		msg = "The requested method is not allowed for the URL"
	case http.StatusBadRequest:
		name = "Bad request"
		msg = "Request header or cookie too large"
	default:
		name = "Internal Server Error"
		msg = "The server encountered an internal error or misconfiguration and was unable to complete your request"
	}

	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}

	errmap := map[string]string{
		"Code": fmt.Sprint(statusCode),
		"Name": name,
		"Msg":  msg,
	}

	if err := t.Execute(w, errmap); err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
}
