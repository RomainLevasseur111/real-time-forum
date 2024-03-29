package research

import (
	"database/sql"
	"errors"

	"real-time-forum/initial"
)

func InsertCategories(cat1, cat2 string) (res []*string, err error) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	if cat1 != "" {
		res = append(res, &cat1)
		row := db.QueryRow("SELECT name FROM CATEGORIES WHERE name = ?;", cat1)
		err = row.Scan()
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO CATEGORIES VALUES(?, 1);", cat1)
		} else {
			_, err = db.Exec("UPDATE CATEGORIES SET posts = posts + 1 WHERE name = ?;", cat1)
		}
		if err != nil {
			return nil, err
		}
	}

	if cat2 != "" {
		res = append(res, &cat2)
		row := db.QueryRow("SELECT name FROM CATEGORIES WHERE name = ?;", cat2)
		err = row.Scan()
		if err == sql.ErrNoRows {
			_, err = db.Exec("INSERT INTO CATEGORIES VALUES(?, 1);", cat2)
		} else {
			_, err = db.Exec("UPDATE CATEGORIES SET posts = posts + 1 WHERE name = ?;", cat2)
		}
		if err != nil {
			return nil, err
		}
	}

	if res == nil {
		err = errors.New("no categories inputed")
		return nil, err
	}

	return res, nil
}

func GetCategories() (cats []initial.CATEGORY) {
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		return nil
	}

	rows, err := db.Query("SELECT name, posts FROM CATEGORIES ORDER BY name;")
	if err != nil {
		return nil
	}

	for rows.Next() {
		var c initial.CATEGORY
		err = rows.Scan(&c.Name, &c.Posts)
		if err != nil {
			return nil
		}
		cats = append(cats, c)
	}

	return
}
