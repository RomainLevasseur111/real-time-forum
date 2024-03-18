package main

import (
	"database/sql"
	"fmt"
	"time"
)

func Publish(userId, category1, category2, content string) int {
	categories, err := InsertCategories(category1, category2)
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	for len(categories) < 2 {
		categories = append(categories, nil)
	}

	res, err := db.Exec(`INSERT INTO "POSTS" ("userid", "category", "categoryB", "content", "postdate") VALUES (?, ?, ?, ?, ?);`,
		userId,
		categories[0],
		categories[1],
		content,
		time.Now().Format(DATEFMT),
	)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	// Get the last inserted ID
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return int(lastInsertID)

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
func GetComments(commentid string) (res []*POST, err error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	sts := `SELECT postid ,userid, category, categoryB, content, postdate FROM POSTS WHERE commentid = ? ORDER BY postdate DESC;`
	rows, err := db.Query(sts, commentid)
	for rows.Next() {
		var post POST
		err = rows.Scan(
			&post.Postid,
			&post.Userid,
			&post.Category,
			&post.CategoryB,
			&post.Content,
			&post.Postdate,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		post.Postdate = FormatDate(post.Postdate)
		post.Likes, post.Dislikes, post.Comments, err = GetStats(post.Postid)
		if err != nil {
			return nil, err
		}

		res = append(res, &post)
	}
	return res, err
}

func FormatDate(value string) (res string) {
	parsedTime, err := time.Parse("2006-01-02T15:04:05Z", value)
	if err != nil {
		return value
	}

	now := time.Now()

	temp := now.Format("2006-01-02T15:04:05Z")
	now, _ = time.Parse("2006-01-02T15:04:05Z", temp)

	diff := now.Sub(parsedTime)

	years := int(diff.Hours()) / 8640
	months := int(diff.Hours()) / 720
	days := int(diff.Hours()) / 24
	hours := int(diff.Hours())
	minutes := int(diff.Minutes())
	seconds := int(diff.Seconds())

	switch {
	case int(years) == 1:
		return "1 year ago"
	case int(years) > 1:
		return fmt.Sprintf("%d years ago", int(years))
	case months == 1:
		return "1 month ago"
	case months > 1:
		return fmt.Sprintf("%d months ago", int(months))
	case days == 1:
		return "1 day ago"
	case days > 1:
		return fmt.Sprintf("%d days ago", int(days))
	case hours == 1:
		return "1 hour ago"
	case hours > 1:
		return fmt.Sprintf("%d hours ago", int(hours))
	case minutes == 1:
		return "1 minute ago"
	case minutes > 1:
		return fmt.Sprintf("%d minutes ago", int(minutes))
	case seconds == 1:
		return "1 second ago"
	case seconds > 1:
		return fmt.Sprintf("%d seconds ago", int(seconds))
	default:
		return "Just now"
	}
}

func GetOnePost(postid string) (*POST, error) {
	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var lol POST
	sts := `SELECT postid , userid, category, categoryB, content, postdate FROM POSTS WHERE postid = ?;`
	row := db.QueryRow(sts, postid)
	err = row.Scan(
		&lol.Postid,
		&lol.Userid,
		&lol.Category,
		&lol.CategoryB,
		&lol.Content,
		&lol.Postdate,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	lol.Likes, lol.Dislikes, lol.Comments, err = GetStats(lol.Postid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	lol.Postdate = FormatDate(lol.Postdate)
	return &lol, err
}

func Comment(postid, userId, category1, category2, content string) int {
	categories, err := InsertCategories(category1, category2)
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open(DRIVER, DB)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	for len(categories) < 2 {
		categories = append(categories, nil)
	}

	res, err := db.Exec(`INSERT INTO "POSTS" ("commentid", "userid", "category", "categoryB", "content", "postdate") VALUES (?,?, ?, ?, ?, ?);`,
		postid,
		userId,
		categories[0],
		categories[1],
		content,
		time.Now().Format(DATEFMT),
	)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	// Get the last inserted ID
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return int(lastInsertID)
}