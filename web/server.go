package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"real-time-forum/initial"
	"real-time-forum/research"
)

func Server(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		initial.Error(w, http.StatusNotFound, "")
		return
	}

	cookie, err := r.Cookie("sessionID")
	data := map[string]interface{}{
		"IP": initial.IP,
	}
	// log the login / register page
	if err != nil {
		tmpl, err := template.ParseFiles("./templates/index.html")
		if err != nil {
			fmt.Print("Error in finding the template")
			fmt.Println(err)
			initial.Error(w, http.StatusInternalServerError, "")
			return
		}
		if tmpl.Execute(w, data) != nil {
			fmt.Print("Error in executing the template")
			fmt.Println(tmpl.Execute(w, data))
			initial.Error(w, http.StatusInternalServerError, "")
			return
		}
		return
	}

	posts, err := GetAllPosts()
	if err != nil {
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}

	users, err := research.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}

	var connectedUser *initial.USER

	if cookie != nil {
		connectedUser, _ = initial.CheckCookie(cookie)
	}

	// log the homepage
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}
	defer db.Close()

	data["ConnectedUser"] = connectedUser
	data["Posts"] = posts
	data["Categories"] = research.GetCategories()
	data["Users"] = users

	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		fmt.Print("Error in finding the template")
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}
	if tmpl.Execute(w, data) != nil {
		fmt.Print("Error in executing the template")
		fmt.Println(tmpl.Execute(w, data))
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}
}
