package main

import "time"

const (
	DB     = "./data.db"
	DRIVER = "sqlite3"
	IP     = "192.168.101.10"
)

type USER struct { // user structure
	NickName, Gender, FirstName, LastName, Email, Pfp string
	Id, Age                                           int
	CreationDate                                      time.Time
}
