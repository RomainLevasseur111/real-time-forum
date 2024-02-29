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
