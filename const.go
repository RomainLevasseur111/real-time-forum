package main

import "time"

const (
	DB     = "./data.db"
	DRIVER = "sqlite3"
	IP     = "192.168.100.249"
)

type USER struct { // user structure
	NickName, Gender, FirstName, LastName, Email, Pfp string
	Id, Age                                           int
	CreationDate                                      time.Time
}
