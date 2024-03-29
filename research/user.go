package research

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"real-time-forum/initial"
)

func GetAllUsers() (users []initial.USER, err error) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT nickname, pfp FROM USERS;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user initial.USER
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

func GetConversation(name1, name2 string) (messages []initial.MESSAGES, err error) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT sendername, receivername, date, pfp, content FROM MESSAGES WHERE sendername = ? AND receivername = ? OR sendername = ? AND receivername = ? ORDER BY messageid ASC;", name1, name2, name2, name1)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message initial.MESSAGES
		err = rows.Scan(
			&message.Sendername,
			&message.Receivername,
			&message.Date,
			&message.Pfp,
			&message.Content,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)

	}

	return messages, nil
}

func UserAction(writer http.ResponseWriter, request *http.Request) {
	referer := request.Referer()
	if request.Method != "POST" {
		initial.Error(writer, http.StatusInternalServerError, "")
		return
	}
	cookie, err := request.Cookie("sessionID")
	if err != nil {
		http.Redirect(writer, request, referer, http.StatusMovedPermanently)
		return
	}
	user, err := initial.CheckCookie(cookie)
	if err != nil {
		fmt.Println(err)
		initial.Error(writer, http.StatusInternalServerError, "")
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
		db, err := sql.Open(initial.DRIVER, initial.DB)
		if err != nil {
			fmt.Println(err)
			initial.Error(writer, http.StatusInternalServerError, "")
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
					initial.Error(writer, http.StatusInternalServerError, "")
					return
				}
			} else {
				sql := `INSERT INTO LikesDislikes (userid, postid, likestatus) VALUES (?, ?, ?);
				`
				_, err = db.Exec(sql, user.Id, postid, action, user.Id)
				if err != nil {
					fmt.Println(err)
					initial.Error(writer, http.StatusInternalServerError, "")
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
						initial.Error(writer, http.StatusInternalServerError, "")
						return
					}
				} else {
					_, err := db.Exec(`DELETE FROM LikesDislikes WHERE userid = ? AND postid = ?;
					`, user.Id, postid)
					if err != nil {
						fmt.Println(err)
						initial.Error(writer, http.StatusInternalServerError, "")
						return
					}
				}
			} else {
				if action == "like" {
					_, err := db.Exec(`UPDATE LikesDislikes SET likestatus = ? WHERE userid = ? AND postid = ?;
					`, action, user.Id, postid)
					if err != nil {
						fmt.Println(err)
						initial.Error(writer, http.StatusInternalServerError, "")
						return
					}
				} else {
					_, err := db.Exec(`UPDATE LikesDislikes SET likestatus = ? WHERE userid = ? AND postid = ?;
					`, action, user.Id, postid)
					if err != nil {
						fmt.Println(err)
						initial.Error(writer, http.StatusInternalServerError, "")
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

func GetOneUser(id string) (*initial.USER, error) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var lol initial.USER
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

func GetOneUserNickname(nickname string) (*initial.USER, error) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var lol initial.USER
	sts := `SELECT nickname, pfp FROM USERS WHERE nickname = ?`
	row := db.QueryRow(sts, nickname)
	err = row.Scan(
		&lol.NickName,
		&lol.Pfp,
	)
	if err != nil {
		return nil, err
	}

	return &lol, err
}

// Sort users by alphabetical order
func SortUsers(users []initial.USER) []initial.USER {
	for i := range users {
		for j := i; j < len(users); j++ {
			if strings.ToLower(users[i].NickName) > strings.ToLower(users[j].NickName) {
				users[i], users[j] = users[j], users[i]
			}
		}
	}
	return users
}

// Check if conversation between two users exists
func ConvExist(loggeduser *initial.USER) ([]initial.USER, error) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	loggedname := loggeduser.NickName + " "
	var messages []initial.MESSAGES
	query := `SELECT messageid, sendername, receivername FROM MESSAGES WHERE sendername = ? OR receivername = ? ORDER BY messageid DESC;`
	rows, err := db.Query(query, loggedname, loggedname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message initial.MESSAGES
		if err := rows.Scan(&message.Messageid, &message.Sendername, &message.Receivername); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	var names []string
	for _, message := range messages {
		if message.Sendername == loggedname {
			dup := false
			for _, i := range names {
				if message.Receivername[:len(message.Receivername)-1] == i {
					dup = true
					break
				}
			}
			if !dup {
				names = append(names, message.Receivername[:len(message.Receivername)-1])
			}
		} else if message.Receivername == loggedname {
			dup := false
			for _, i := range names {
				if message.Sendername[:len(message.Sendername)-1] == i {
					dup = true
					break
				}
			}
			if !dup {
				names = append(names, message.Sendername[:len(message.Sendername)-1])
			}
		}
	}
	var result []initial.USER
	for _, name := range names {
		user, err := GetOneUserNickname(name)
		if err != nil {
			fmt.Println(err)
			break
		}
		result = append(result, *user)
	}

	return result, nil
}

//
