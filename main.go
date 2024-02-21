package main

import (
	"fmt"
	"net/http"
)

func main() {
	CreateDatabase()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server starting...\nServer hosted at http://localhost:8989/")

	http.HandleFunc("/", Server)
	if err := http.ListenAndServe(":8989", nil); err != nil {
		fmt.Println("Error starting the server")
		return
	}
}
