package main

import "database/sql"

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
