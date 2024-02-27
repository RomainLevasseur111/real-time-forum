package main

import (
	"fmt"
	"net/http"
)

func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients = append(clients, conn)

	defer func() {
		// Remove the connection from the clients slice when it's closed
		for i, client := range clients {
			if client == conn {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		for _, client := range clients {
			if err = client.WriteMessage(msgType, msg); err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}
