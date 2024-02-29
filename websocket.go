package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func Chat_Websocket(w http.ResponseWriter, r *http.Request) {
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

		// Remove the connection from the clients in the struct
		for i, c := range connection {
			if c.Conn == conn {
				connection = append(connection[:i], connection[i+1:]...)
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

		if !strings.Contains(string(msg), " ") {

			var conn_ CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg)

			fmt.Printf("Connection of %s from %s\n", conn_.Name, conn.RemoteAddr())

			connection = append(connection, conn_)
			continue
		}

		// format: sendername/receivername/date/content
		msgData := strings.SplitAfterN(string(msg), " ", 4)

		if len(msgData) != 4 || msgData[3] == "" {
			fmt.Println("empty message")
			continue
		}

		var conn_ CONNECTIONS
		conn_.Conn = conn
		conn_.Name = msgData[0][:len(msgData[0])-1]

		fmt.Printf("%s\n", string(msg))

		connection = append(connection, conn_)

		db, err := sql.Open(DRIVER, DB)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()

		rows := db.QueryRow("SELECT pfp FROM USERS WHERE nickname = ?", msgData[0][:len(msgData[0])-1])
		var pfp string
		err = rows.Scan(&pfp)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = db.Exec(`INSERT INTO MESSAGES VALUES(NULL, ?, ?, ?, ?, ?)`,
			msgData[0],
			msgData[1],
			msgData[2],
			pfp,
			msgData[3])

		if err != nil {
			fmt.Println(err)
			return
		}

		temp := msgData[0] + " " + msgData[2] + " " + pfp + " " + msgData[3] + " "

		count := 0
		for _, client := range clients {
			for _, c := range connection {
				if (c.Name == msgData[1][:len(msgData[1])-1] || c.Name == msgData[0][:len(msgData[0])-1]) && c.Conn == client {
					count++
					if err = client.WriteMessage(msgType, []byte(temp)); err != nil {
						fmt.Println(err)
						break
					}
					break
				}
			}
			if count == 2 {
				break
			}
		}
	}
}
