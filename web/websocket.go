package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"real-time-forum/initial"
	"real-time-forum/research"
)

func Chat_Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := initial.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Remove the connection of a client when it's closed
	defer func() {
		for i, c := range initial.Chat_connection {
			if c.Conn == conn {
				initial.Chat_connection = append(initial.Chat_connection[:i], initial.Chat_connection[i+1:]...)
				break
			}
		}
	}()

	// Send a private message
	send := func(msgData []string, pfp, name1, name2 string, msgType int) {
		temp := msgData[0] + " " + msgData[2] + " " + pfp + " " + msgData[3] + " "

		for _, c := range initial.Chat_connection {
			if c.Name == name1 || c.Name == name2 {
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
			for _, c := range initial.Chat_connection {
				if c.Name == strings.Split(string(msg), " ")[1] {
					response = "IsCo_Yes"
				}
			}

			for _, c := range initial.Chat_connection {
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

			var conn_ initial.CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg)

			fmt.Printf("Connection of %s from %s at the chat_websocket\n", conn_.Name, conn.RemoteAddr())

			initial.Chat_connection = append(initial.Chat_connection, conn_)
			continue
		}

		// Ask for all the client in the database
		if string(msg)[0:4] == "U_N " {

			loggedUser, err := research.GetOneUserNickname(string(msg[4:]))
			if err != nil {
				fmt.Println(err)
				return
			}

			users, err := research.GetAllUsers()
			if err != nil {
				fmt.Println(err)
				return
			}
			temp := "U_N "

			convusers, err := research.ConvExist(loggedUser)
			if err != nil {
				fmt.Println(err)
				break
			}
			research.SortUsers(users)
			for _, i := range users {
				isIn := false
				for _, j := range convusers {
					if i == j {
						isIn = true
						break
					}
				}
				if !isIn {
					convusers = append(convusers, i)
				}
			}
			for _, user := range convusers {
				isConnected := "../static/img/disconnected.webp"
				for _, c := range initial.Chat_connection {
					if c.Name == user.NickName {
						isConnected = "../static/img/connected.png"
					}
				}
				temp += user.NickName + " " + user.Pfp + " " + isConnected + " "
			}

			for _, c := range initial.Chat_connection {
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
			messages, err := research.GetConversation(msgData[1], msgData[2])
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, message := range messages {
				var temp []string
				temp = append(temp, message.Sendername, message.Receivername, message.Date, message.Content)
				send(temp, message.Pfp, msgData[2][:len(msgData[2])-1], "", msgType)
			}

			continue
		}

		var conn_ initial.CONNECTIONS
		conn_.Conn = conn
		conn_.Name = msgData[0][:len(msgData[0])-1]

		fmt.Printf("%s\n", string(msg))

		db, err := sql.Open(initial.DRIVER, initial.DB)
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

		send(msgData, pfp, msgData[1][:len(msgData[1])-1], msgData[0][:len(msgData[0])-1], msgType)
	}
}

func Post_Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := initial.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Remove the connection of a client when it's closed
	defer func() {
		for i, c := range initial.Post_connection {
			if c.Conn == conn {
				initial.Post_connection = append(initial.Post_connection[:i], initial.Post_connection[i+1:]...)
				break
			}
		}
	}()

	// Display a given post
	displayPost := func(pfp, nickname, content string, category, categoryB *string, msgType, postid int, userNameToSend string) {
		cat1, cat2 := "_&nbsp_", "_&nbsp_"
		if category != nil {
			cat1 = *category
		}
		if categoryB != nil {
			cat2 = *categoryB
		}
		temp := "P_B " + pfp + " " + nickname + " " + cat1 + " " + cat2 + " " + strconv.Itoa(postid) + " " + content
		for _, c := range initial.Post_connection {
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

			var conn_ initial.CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg)

			fmt.Printf("Connection of %s from %s at the post_websocket\n", conn_.Name, conn.RemoteAddr())

			initial.Post_connection = append(initial.Post_connection, conn_)
			continue
		}

		// Publish a post to all user
		if string(msg[0:4]) == "P_B " {
			msgData := strings.SplitN(string(msg), " ", 5)
			if msgData[4] == "" {
				fmt.Println("Empty Post")
				continue
			}
			postid := Publish(msgData[1], msgData[2], msgData[3], msgData[4])
			user, err := research.GetOneUser(msgData[1])
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
			var conn_ initial.CONNECTIONS
			conn_.Conn = conn
			conn_.Name = string(msg[4:])

			fmt.Printf("Connection of %s from %s at the post_websocket\n", conn_.Name, conn.RemoteAddr())

			initial.Post_connection = append(initial.Post_connection, conn_)

			// Display all post
			userNameToSend := string(msg[4:])

			posts, err := GetAllPosts()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, post := range posts {
				user, err := research.GetOneUser(strconv.Itoa(post.Userid))
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
	conn, err := initial.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	posttype := ""
	defer conn.Close()
	// Remove the connection of a client when it's closed
	defer func() {
		for i, c := range initial.Comment_connection {
			if c.Conn == conn {
				initial.Comment_connection = append(initial.Comment_connection[:i], initial.Comment_connection[i+1:]...)
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

		displayComment := func(name, originalPostId, posttype, pfp, nickname, content string, category, categoryB *string, msgType, postid int) {
			cat1, cat2 := "_&nbsp_", "_&nbsp_"
			if category != nil {
				cat1 = *category
			}
			if categoryB != nil {
				cat2 = *categoryB
			}
			temp := posttype + pfp + " " + nickname + " " + cat1 + " " + cat2 + " " + strconv.Itoa(postid) + " " + content
			for _, c := range initial.Comment_connection {
				if c.Name == name || (name == "ALL" && originalPostId == c.CommentId) {
					if err = c.Conn.WriteMessage(msgType, []byte(temp)); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}

		// Show all comments of a post to a new user
		if string(msg[0:4]) == "C_M " {

			// Register new connection
			var conn_ initial.CONNECTIONS
			conn_.Conn = conn
			conn_.Name = strings.Split(string(msg), " ")[1]
			conn_.CommentId = strings.Split(string(msg), " ")[2]

			fmt.Printf("Connection of %s from %s at the comment_websocket\n", conn_.Name, conn.RemoteAddr())

			initial.Comment_connection = append(initial.Comment_connection, conn_)

			post, err := GetOnePost(conn_.CommentId)
			if err != nil {
				fmt.Println(err)
				break
			}
			postuser, err := research.GetOneUser(strconv.Itoa(post.Userid))
			if err != nil {
				fmt.Println(err)
				break
			}
			// Initial post
			posttype = "P_M "
			displayComment(conn_.Name, "", posttype, postuser.Pfp, postuser.NickName, post.Content, post.Category, post.CategoryB, msgType, post.Postid)

			comments, err := GetComments(conn_.CommentId)
			if err != nil {
				fmt.Println(err)
				break
			}

			for _, comment := range comments {
				user, err := research.GetOneUser(strconv.Itoa(comment.Userid))
				if err != nil {
					fmt.Println(err)
					break
				}
				// All its comments
				posttype = "C_M "
				displayComment(conn_.Name, "", posttype, user.Pfp, user.NickName, comment.Content, comment.Category, comment.CategoryB, msgType, comment.Postid)
			}

			// New comment
		} else if string(msg[0:4]) == "P_C " {
			msgData := strings.SplitN(string(msg), " ", 4)
			if msgData[3] == "" {
				fmt.Println("Empty comment")
				continue
			}

			newCommentId := Comment(msgData[1], msgData[2], msgData[3])
			post, err := GetOnePost(strconv.Itoa(newCommentId))
			if err != nil {
				fmt.Println(err)
				break
			}

			postuser, err := research.GetOneUser(strconv.Itoa(post.Userid))
			if err != nil {
				fmt.Println(err)
				break
			}
			displayComment("ALL", msgData[1], "C_M ", postuser.Pfp, postuser.NickName, post.Content, nil, nil, msgType, newCommentId)
		}
	}
}
