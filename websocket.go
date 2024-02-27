package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
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

		// format: sendername/receivername/date/content
		msgData := strings.SplitAfterN(string(msg), " ", 4)

		db, err := sql.Open(DRIVER, DB)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()

		_, err = db.Exec(`INSERT INTO MESSAGES VALUES(NULL, ?, ?, ?, ?)`,
			msgData[0],
			msgData[1],
			msgData[2],
			msgData[3])

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, client := range clients {
			if err = client.WriteMessage(msgType, msg); err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}
