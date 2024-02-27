package main

import (
	"fmt"
	"net/http"
)

func main() {
	CreateDatabase()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server starting...\nServer hosted at http://" + IP + ":8989/")

	http.HandleFunc("/", Server)
	http.HandleFunc("/register", Registration)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/echo", Websocket)

	if err := http.ListenAndServe("0.0.0.0:8989", nil); err != nil {
		fmt.Println("Error starting the server")
		return
	}
}
