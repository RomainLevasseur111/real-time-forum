package main

import (
	"fmt"
	"net/http"

	"real-time-forum/initial"
	"real-time-forum/login"
	"real-time-forum/research"
	"real-time-forum/web"
)

func init() {
	initial.CreateDatabase()
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server starting...\nServer hosted at http://" + initial.IP + ":8989/")

	http.HandleFunc("/", web.Server)
	http.HandleFunc("/register", login.Registration)
	http.HandleFunc("/login", login.Login)
	http.HandleFunc("/logout", login.Logout)
	http.HandleFunc("/chat_websocket", web.Chat_Websocket)
	http.HandleFunc("/useraction", research.UserAction)
	http.HandleFunc("/post_websocket", web.Post_Websocket)
	http.HandleFunc("/comment_websocket", web.Comment_Websocket)

	if err := http.ListenAndServe("0.0.0.0:8989", nil); err != nil {
		fmt.Println("Error starting the server")
		fmt.Println(err)
		return
	}
}
