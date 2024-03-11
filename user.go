package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func GetAllUsers() (users []USER, err error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT nickname, pfp FROM USERS;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user USER
		err = rows.Scan(
			&user.NickName,
			&user.Pfp,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}

func GetConversation(name1, name2 string) (messages []MESSAGES, err error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT sendername, receivername, date, pfp, content FROM MESSAGES WHERE sendername = ? AND receivername = ?;", name1, name2)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message MESSAGES
		err = rows.Scan(
			&message.sendername,
			&message.receivername,
			&message.date,
			&message.pfp,
			&message.content,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)

	}

	rows_, err := db.Query("SELECT sendername, receivername, date, pfp, content FROM MESSAGES WHERE sendername = ? AND receivername = ?;", name2, name1)
	if err != nil {
		return nil, err
	}

	for rows_.Next() {
		var message MESSAGES
		err = rows_.Scan(
			&message.sendername,
			&message.receivername,
			&message.date,
			&message.pfp,
			&message.content,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)

	}

	return messages, nil
}
func userAction(writer http.ResponseWriter, request *http.Request) {
	referer := request.Referer()
	if request.Method != "POST" {
		Error(writer, http.StatusInternalServerError)
		return
	}
	cookie, err := request.Cookie("sessionID")
	if err != nil {
		http.Redirect(writer, request, referer, http.StatusMovedPermanently)
		return
	}
	user, err := checkCookie(cookie)
	if err != nil {
		fmt.Println(err)
		Error(writer, http.StatusInternalServerError)
		return
	}

	action := request.FormValue("action")

	postid := request.FormValue("postid")
	if postid == "" {
		http.Redirect(writer, request, referer, http.StatusMovedPermanently)
		return
	}
	switch action {
	case "like", "dislike":
		db, err := sql.Open(DRIVER, DB)
		if err != nil {
			fmt.Println(err)
			Error(writer, http.StatusInternalServerError)
			return
		}

		row := db.QueryRow("SELECT likestatus FROM LikesDislikes WHERE userid = ? AND postid = ?;", user.Id, postid)
		var res string
		err = row.Scan(&res)
		if err == sql.ErrNoRows {
			if action == "like" {
				sql := `INSERT INTO LikesDislikes (userid, postid, likestatus) VALUES (?, ?, ?);
				`
				_, err = db.Exec(sql, user.Id, postid, action, user.Id)
				if err != nil {
					fmt.Println(err)
					Error(writer, http.StatusInternalServerError)
					return
				}
			} else {
				sql := `INSERT INTO LikesDislikes (userid, postid, likestatus) VALUES (?, ?, ?);
				`
				_, err = db.Exec(sql, user.Id, postid, action, user.Id)
				if err != nil {
					fmt.Println(err)
					Error(writer, http.StatusInternalServerError)
					return
				}
			}
		} else {
			if action == res {
				if action == "like" {
					_, err := db.Exec(`DELETE FROM LikesDislikes WHERE userid = ? AND postid = ?;
					`, user.Id, postid, user.Id)
					if err != nil {
						fmt.Println(err)
						Error(writer, http.StatusInternalServerError)
						return
					}
				} else {
					_, err := db.Exec(`DELETE FROM LikesDislikes WHERE userid = ? AND postid = ?;
					`, user.Id, postid)
					if err != nil {
						fmt.Println(err)
						Error(writer, http.StatusInternalServerError)
						return
					}
				}
			} else {
				if action == "like" {
					_, err := db.Exec(`UPDATE LikesDislikes SET likestatus = ? WHERE userid = ? AND postid = ?;
					`, action, user.Id, postid)
					if err != nil {
						fmt.Println(err)
						Error(writer, http.StatusInternalServerError)
						return
					}
				} else {
					_, err := db.Exec(`UPDATE LikesDislikes SET likestatus = ? WHERE userid = ? AND postid = ?;
					`, action, user.Id, postid)
					if err != nil {
						fmt.Println(err)
						Error(writer, http.StatusInternalServerError)
						return
					}
				}
			}
		}
		http.Redirect(writer, request, "/", http.StatusMovedPermanently)

	default:
		http.Redirect(writer, request, "/", http.StatusMovedPermanently)
	}
}
func GetOneUser(id string) (*USER, error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var lol USER
	sts := `SELECT nickname, pfp FROM USERS WHERE id = ?`
	row := db.QueryRow(sts, id)
	err = row.Scan(
		&lol.NickName,
		&lol.Pfp,
)
	if err != nil {
		return nil, err
	}

	return &lol, err
}