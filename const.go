package main

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	DB     = "./data.db"
	DRIVER = "sqlite3"
	IP     = "192.168.100.250"
)

var (
	clients  []*websocket.Conn
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type USER struct { // user structure
	NickName, Gender, FirstName, LastName, Email, Pfp string
	Id, Age                                           int
	CreationDate                                      time.Time
}
