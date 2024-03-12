package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Chat_Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	displayed := false

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

	send := func(msgData []string, pfp string, msgType int) {
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

	displayPost := func(pfp, nickname, content string, category, categoryB *string, msgType int) {
		cat1, cat2 := "_&nbsp_", "_&nbsp_"
		if category != nil {
			cat1 = *category
		}
		if categoryB != nil {
			cat2 = *categoryB
		}
		temp := "PUBLISH_ " + pfp + " " + nickname + " " + cat1 + " " + cat2 + " " + content
		for _, client := range clients {
			if err = client.WriteMessage(msgType, []byte(temp)); err != nil {
				fmt.Println(err)
				break
			}
		}
	}
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		if string(msg[0:8]) == "PUBLISH_" {
			msgData := strings.SplitN(string(msg), " ", 5)
			Publish(msgData[1], msgData[2], msgData[3], msgData[4])
			user, err := GetOneUser(msgData[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			displayPost(user.Pfp, user.NickName, msgData[4], &msgData[2], &msgData[3], msgType)

		} else {

			if !strings.Contains(string(msg), " ") {

				var conn_ CONNECTIONS
				conn_.Conn = conn
				conn_.Name = string(msg)

				fmt.Printf("Connection of %s from %s\n", conn_.Name, conn.RemoteAddr())

				connection = append(connection, conn_)
				continue
			}

			if string(msg)[0:4] == "U_N " {
				users, err := GetAllUsers()
				if err != nil {
					fmt.Println(err)
					return
				}
				if !displayed {
					posts, err := GetAllPosts()
					if err != nil {
						fmt.Println(err)
						return
					}
					for _, post := range posts {
						user, err := GetOneUser(strconv.Itoa(post.Userid))
						if err != nil {
							fmt.Println(err)
							return
						}
						displayPost(user.Pfp, user.NickName, post.Content, post.Category, post.CategoryB, msgType)
					}
					displayed = true
				}

				// sort users by alphabetical order
				for i := range users {
					for j := i; j < len(users); j++ {
						if strings.ToLower(users[i].NickName) > strings.ToLower(users[j].NickName) {
							users[i], users[j] = users[j], users[i]
						}
					}
				}

				temp := "U_N "
				for _, user := range users {
					isConnected := "../static/img/disconnected.webp"
					for _, c := range connection {
						if c.Name == user.NickName {
							isConnected = "../static/img/connected.png"
						}
					}
					temp += user.NickName + " " + user.Pfp + " " + isConnected + " "
				}

				for _, client := range clients {
					for _, c := range connection {
						if string(msg)[4:] == c.Name && c.Conn == client {
							if err = client.WriteMessage(msgType, []byte(temp)); err != nil {
								fmt.Println(err)
								break
							}
							temp = ""
							break
						}
					}
					if temp == "" {
						break
					}
				}

				continue
			}

			// format: sendername/receivername/date/content
			msgData := strings.SplitAfterN(string(msg), " ", 4)

			if len(msgData) != 4 || msgData[3] == "" {
				fmt.Println("empty message")
				continue
			}

			if msgData[0] == "GAM " && msgData[3] == "_" {
				messages, err := GetConversation(msgData[1], msgData[2])
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, message := range messages {
					var temp []string
					temp = append(temp, message.sendername, message.receivername, message.date, message.content)
					send(temp, message.pfp, msgType)
				}

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

			send(msgData, pfp, msgType)
		}

	}
}
