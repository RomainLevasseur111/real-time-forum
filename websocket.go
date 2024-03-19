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

	// Remove the connection of a client when it's closed
	defer func() {
		for i, c := range chat_connection {
			if c.Conn == conn {
				chat_connection = append(chat_connection[:i], chat_connection[i+1:]...)
				break
			}
		}
	}()

	// Send a private message
	send := func(msgData []string, pfp string, msgType int) {
		temp := msgData[0] + " " + msgData[2] + " " + pfp + " " + msgData[3] + " "

		for _, c := range chat_connection {
			if c.Name == msgData[1][:len(msgData[1])-1] || c.Name == msgData[0][:len(msgData[0])-1] {
				if err = c.Conn.WriteMessage(msgType, []byte(temp)); err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		// Is the client given connected ?
		if string(msg[0:5]) == "IsCo " && len(strings.Split(string(msg), " ")) == 3 {
			response := "IsCo_No"
			for _, c := range chat_connection {
				if c.Name == strings.Split(string(msg), " ")[1] {
					response = "IsCo_Yes"
				}
			}

			for _, c := range chat_connection {
				if c.Name == strings.Split(string(msg), " ")[2] {
					if err = c.Conn.WriteMessage(msgType, []byte(response)); err != nil {
						fmt.Println(err)
						break
					}
					break
				}
			}

			continue
		}

		// New connection of a client
		if !strings.Contains(string(msg), " ") {

			var conn_ CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg)

			fmt.Printf("Connection of %s from %s at the chat_websocket\n", conn_.Name, conn.RemoteAddr())

			chat_connection = append(chat_connection, conn_)
			continue
		}

		// Ask for all the client in the database
		if string(msg)[0:4] == "U_N " {
			users, err := GetAllUsers()
			if err != nil {
				fmt.Println(err)
				return
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
				for _, c := range chat_connection {
					if c.Name == user.NickName {
						isConnected = "../static/img/connected.png"
					}
				}
				temp += user.NickName + " " + user.Pfp + " " + isConnected + " "
			}

			for _, c := range chat_connection {
				if string(msg)[4:] == c.Name {
					if err = c.Conn.WriteMessage(msgType, []byte(temp)); err != nil {
						fmt.Println(err)
						break
					}
					temp = ""
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

		// Ask for a conversation
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

func Post_Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Remove the connection of a client when it's closed
	defer func() {
		for i, c := range post_connection {
			if c.Conn == conn {
				post_connection = append(post_connection[:i], post_connection[i+1:]...)
				break
			}
		}
	}()

	// Display a given post
	displayPost := func(pfp, nickname, content string, category, categoryB *string, msgType, postid int, userNameToSend string) {
		fmt.Println(userNameToSend)
		fmt.Println(post_connection)
		cat1, cat2 := "_&nbsp_", "_&nbsp_"
		if category != nil {
			cat1 = *category
		}
		if categoryB != nil {
			cat2 = *categoryB
		}
		temp := "P_B " + pfp + " " + nickname + " " + cat1 + " " + cat2 + " " + strconv.Itoa(postid) + " " + content
		for _, c := range post_connection {
			if userNameToSend == "" {
				if err = c.Conn.WriteMessage(msgType, []byte(temp)); err != nil {
					fmt.Println(err)
					return
				}
				continue
			}
			if c.Name == userNameToSend {
				if err = c.Conn.WriteMessage(msgType, []byte(temp)); err != nil {
					fmt.Println(err)
					return
				}
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

		// New connection of a client
		if !strings.Contains(string(msg), " ") {

			var conn_ CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg)

			fmt.Printf("Connection of %s from %s at the post_websocket\n", conn_.Name, conn.RemoteAddr())

			post_connection = append(post_connection, conn_)
			continue
		}

		// Publish a post to all user
		if string(msg[0:4]) == "P_B " {
			msgData := strings.SplitN(string(msg), " ", 5)
			fmt.Println(len(msgData))
			postid := Publish(msgData[1], msgData[2], msgData[3], msgData[4])
			user, err := GetOneUser(msgData[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			displayPost(user.Pfp, user.NickName, msgData[4], &msgData[2], &msgData[3], msgType, postid, "")
			continue
		}

		// Display all post to a new user
		if string(msg[0:4]) == "1_D " {

			// Register new connection
			var conn_ CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg[4:])

			fmt.Printf("Connection of %s from %s at the post_websocket\n", conn_.Name, conn.RemoteAddr())

			post_connection = append(post_connection, conn_)

			// Display all post
			userNameToSend := string(msg[4:])

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
				displayPost(user.Pfp, user.NickName, post.Content, post.Category, post.CategoryB, msgType, post.Postid, userNameToSend)
			}
		}
	}
}

func Comment_Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	posttype := ""
	defer conn.Close()
	// Remove the connection of a client when it's closed
	defer func() {
		for i, c := range comment_connection {
			if c.Conn == conn {
				comment_connection = append(comment_connection[:i], comment_connection[i+1:]...)
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
		displayComment := func(posttype, pfp, nickname, content string, category, categoryB *string, msgType, postid int) {
			cat1, cat2 := "_&nbsp_", "_&nbsp_"
			if category != nil {
				cat1 = *category
			}
			if categoryB != nil {
				cat2 = *categoryB
			}
			temp := posttype + pfp + " " + nickname + " " + cat1 + " " + cat2 + " " + strconv.Itoa(postid) + " " + content
			for _, c := range comment_connection {
				if err = c.Conn.WriteMessage(msgType, []byte(temp)); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		// Show all comments of a post to a new user
		if string(msg[0:4]) == "C_M " {

			// Register new connection
			var conn_ CONNECTIONS
			conn_.Conn = conn
			conn_.Name = strings.Split(string(msg), " ")[1]

			fmt.Printf("Connection of %s from %s at the comment_websocket\n", conn_.Name, conn.RemoteAddr())

			comment_connection = append(comment_connection, conn_)

			post, err := GetOnePost(strings.Split(string(msg), " ")[2])
			if err != nil {
				fmt.Println(err)
				break
			}
			postuser, err := GetOneUser(strconv.Itoa(post.Userid))
			if err != nil {
				fmt.Println(err)
				break
			}
			// Post initial
			posttype = "P_M "
			displayComment(posttype, postuser.Pfp, postuser.NickName, post.Content, post.Category, post.CategoryB, msgType, post.Postid)

			comments, err := GetComments(strings.Split(string(msg), " ")[2])
			if err != nil {
				fmt.Println(err)
				break
			}
			for _, comment := range comments {
				user, err := GetOneUser(strconv.Itoa(comment.Userid))
				if err != nil {
					fmt.Println(err)
					break
				}
				// Tous ses commentaires
				posttype = "C_M "
				displayComment(posttype, user.Pfp, user.NickName, comment.Content, comment.Category, comment.CategoryB, msgType, comment.Postid)
			}
			// PB: a chaque fois qu'on clique sur le bouton comment, ouvre un nouvelle connexion websocket, donc tout se fait autant de fois qu'il y a
			// connexion donc duplique tous les commentaires.
			// conn.Close() permet d'annuler ça, mais je pense pas que ce soit une bonne solution.
		} else if string(msg[0:4]) == "P_C " {
			msgData := strings.SplitN(string(msg), " ", 4)
			newCommentId := Comment(msgData[1], msgData[2], msgData[3])
			post, err := GetOnePost(strconv.Itoa(newCommentId))
			if err != nil {
				fmt.Println(err)
				break
			}
			postuser, err := GetOneUser(strconv.Itoa(post.Userid))
			if err != nil {
				fmt.Println(err)
				break
			}
			displayComment("C_M ",postuser.Pfp, postuser.NickName, post.Content, nil, nil, msgType, newCommentId )
		}
	}
}
