package main

import (
	"database/sql"
	"fmt"
	"time"
)

func Publish(userId, category1, category2, content string) {
	categories, err := InsertCategories(category1, category2)
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	for len(categories) < 2 {
		categories = append(categories, nil)
	}

	// ajouter publication de commentaires : il faut recup l'id du post initial
	_, err = db.Exec(`INSERT INTO "POSTS" ("userid", "category", "categoryB", "content", "postdate") VALUES (?, ?,?, ?, ?);`,
		userId,
		// toInput,
		categories[0],
		categories[1],
		content,
		time.Now().Format(DATEFMT),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//sql := ""
	/*if toInput == nil {
		sql = "UPDATE USERS SET posts = posts + 1 WHERE cookie = ?;"
	} else {
		sql = "UPDATE USERS SET comments = comments + 1 WHERE cookie = ?;"
	}

	_, err = db.Exec(sql, cookie.Value)
	if err != nil {
		fmt.Println(err)
		return
	}*/

}

func GetAllPosts() (posts []POST, err error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT userid, postid, category, categoryB, content FROM POSTS WHERE commentid IS NULL ORDER BY postdate ASC;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post POST
		err = rows.Scan(
			&post.Userid,
			&post.Postid,
			&post.Category,
			&post.CategoryB,
			&post.Content,
		)
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
