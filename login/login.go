package login

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"real-time-forum/initial"
	"real-time-forum/research"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Check login datas in the database
	if r.Method != "POST" {
		initial.Error(w, http.StatusMethodNotAllowed, "")
		return
	}
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}
	defer db.Close()
	var pass string
	name := r.FormValue("login")
	pwd := r.FormValue("password")
	if !research.EmailAlreadyExist(name) {
		rows, err := db.Query("SELECT nickname FROM USERS WHERE email = ?", name)
		if err != nil {
			fmt.Println(err)
			initial.Error(w, http.StatusInternalServerError, "")
			return
		}
		for rows.Next() {
			rows.Scan(&name)
		}
	}
	if !research.NicknameAlreadyExists(name) {
		rows, err := db.Query("SELECT password FROM USERS WHERE nickname = ?", name)
		if err != nil {
			fmt.Println(err)
			initial.Error(w, http.StatusInternalServerError, "")
			return
		}
		for rows.Next() {
			rows.Scan(&pass)
		}
	}
	if research.CheckPasswordHash(pwd, pass) {
		initial.GiveCookie(w, name)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		var erreur string = "We did not recognize your nickname, email or password in our database, you better try again otherwise we will find you and we will kill you."
		errorMsg := map[string]interface{}{
			"Error": erreur,
		}
		tmpl.Execute(w, errorMsg)
	}
}
