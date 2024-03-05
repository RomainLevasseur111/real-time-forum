package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func Publish(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		Error(writer, http.StatusNotFound)
		return
	}

	cookie, err := request.Cookie("sessionID")
	if err != nil {
		fmt.Println(err)
		Error(writer, http.StatusInternalServerError)
		return
	}

	user, err := checkCookie(cookie)
	if err != nil {
		fmt.Println(err)
		Error(writer, http.StatusInternalServerError)
		return
	}

	category1 := request.FormValue("category1")
	category2 := request.FormValue("category2")
	categories, err := InsertCategories(category1, category2)
	if err != nil {
		fmt.Println(err)
	}

	content := request.FormValue("content")
	if content == "" {
		http.Redirect(writer, request, "/?error=grosseerreurmonamimaisjesaispaspourquoijailaflemmedereflechirladessus", http.StatusMovedPermanently)
		return
	}

	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		Error(writer, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	for len(categories) < 2 {
		categories = append(categories, nil)
	}

	// ajouter publication de commentaires : il faut recup l'id du post initial
	_, err = db.Exec(`INSERT INTO "POSTS" ("userid","username", "category", "categoryB","userpfp", "content", "postdate") VALUES (?, ?, ?, ?, ?, ?, ?);`,
		user.Id,
		// toInput,
		user.NickName,
		categories[0],
		categories[1],
		user.Pfp,
		content,
		time.Now().Format(DATEFMT),
	)
	if err != nil {
		fmt.Println(err)
		Error(writer, http.StatusInternalServerError)
		return
	}

	sql := ""
	/*if toInput == nil {
		sql = "UPDATE USERS SET posts = posts + 1 WHERE cookie = ?;"
	} else {
		sql = "UPDATE USERS SET comments = comments + 1 WHERE cookie = ?;"
	}*/

	_, err = db.Exec(sql, cookie.Value)
	if err != nil {
		fmt.Println(err)
		Error(writer, http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, "/", http.StatusMovedPermanently)
}

func GetAllPosts() (posts []POST, err error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT userid, postid, username, category, categoryB, userpfp, content FROM POSTS WHERE commentid IS NULL ORDER BY postdate DESC;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post POST
		err = rows.Scan(
			&post.Userid,
			&post.Postid,
			&post.Username,
			&post.Category,
			&post.CategoryB,
			&post.Userpfp,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, post.Comments, err = GetStats(post.Postid)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil
}

func GetStats(postid int) (likes int, dislikes int, comments int, err error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT COUNT(*) FROM LikesDislikes  WHERE postid = ? AND likestatus = 'like';", postid)
	err = row.Scan(&likes)
	if err != nil {
		return
	}

	row = db.QueryRow("SELECT COUNT(*) FROM LikesDislikes  WHERE postid = ? AND likestatus = 'dislike';", postid)
	err = row.Scan(&dislikes)
	if err != nil {
		return
	}

	row = db.QueryRow("SELECT COUNT(*) FROM POSTS WHERE commentid = ?;", postid)
	err = row.Scan(&comments)
	if err != nil {
		return
	}

	return
}
