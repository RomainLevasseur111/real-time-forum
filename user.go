package main

import (
	"database/sql"
	"fmt"
	"net/http"
)
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