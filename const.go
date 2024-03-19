package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	DB      = "./data.db"
	DRIVER  = "sqlite3"
	IP      = "10.10.1.199"
	DATEFMT = "2006-01-02 15:04:05"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	chat_connection, post_connection, comment_connection []CONNECTIONS
)

type USER struct { // user structure
	NickName, Gender, FirstName, LastName, Email, Pfp string
	Id, Age                                           int
	CreationDate                                      time.Time
}
type POST struct { // post structure
	Postdate, Content                         string
	Category, CategoryB                       *string
	Userid, Postid, Likes, Comments, Dislikes int
}
type CATEGORY struct { // Category structure
	Name  string
	Posts int
}

type CONNECTIONS struct {
	Conn *websocket.Conn
	Name string
}

type MESSAGES struct {
	sendername, receivername, date, pfp, content string
}
