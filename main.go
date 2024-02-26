package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients  []*websocket.Conn
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	CreateDatabase()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Server starting...\nServer hosted at http://" + IP + ":8989/")

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		clients = append(clients, conn)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range clients {
				if err = client.WriteMessage(msgType, msg); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	})

	http.HandleFunc("/", Server)
	http.HandleFunc("/register", Registration)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)

	if err := http.ListenAndServe("0.0.0.0:8989", nil); err != nil {
		fmt.Println("Error starting the server")
		return
	}
}
