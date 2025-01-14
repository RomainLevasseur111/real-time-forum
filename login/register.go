package login

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"real-time-forum/initial"
	"real-time-forum/research"

	"golang.org/x/crypto/bcrypt"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		initial.Error(w, http.StatusMethodNotAllowed, "")
		return
	}

	// Check if registeration datas are correct
	nickname := r.FormValue("nickname")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	firstname := r.FormValue("firstName")
	lastname := r.FormValue("lastName")
	email := r.FormValue("email")
	password := r.FormValue("password")

	var errMsg string
	if len(nickname) > 9 {
		errMsg = "Nickname can't have more than 9 characters"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}
	errMsg = CheckName(nickname)
	if errMsg != "" {
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}
	if !research.NicknameAlreadyExists(nickname) {
		errMsg = "Nickname already exists"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	errMsg = CheckName(firstname)
	if errMsg != "" {
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	errMsg = CheckName(lastname)
	if errMsg != "" {
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		errMsg = "Invalid age format"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}
	if ageInt < 0 {
		errMsg = "Age can't be negative"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	if gender != "male" && gender != "female" && gender != "other" {
		errMsg = "Invalid gender format"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		errMsg = "Invalid email address"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}
	if !research.EmailAlreadyExist(email) {
		errMsg = "Email already exists"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	if len(password) > 30 {
		errMsg = "Password can't have more than 30 characters"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}
	if len(password) < 8 {
		errMsg = "Password can't have less than 8 characters"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}
	psw, err := HashPassword(password)
	if err != nil {
		errMsg = "Error hashing your password, try another password"
		fmt.Println(errMsg)
		initial.Error(w, http.StatusNotAcceptable, errMsg)
		return
	}

	// If registration is successful:
	// Store datas in database
	db, err := sql.Open(initial.DRIVER, initial.DB)
	if err != nil {
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}
	defer db.Close()

	today := time.Now()
	_, err = db.Exec(`INSERT INTO USERS VALUES(NULL, ?, ?, ?, ?, ?, ?, ?, "https://i.stack.imgur.com/l60Hf.png", ?, "", NULL)`,
		nickname,
		ageInt,
		gender,
		firstname,
		lastname,
		email,
		psw,
		today)

	if err != nil {
		fmt.Println(err)
		initial.Error(w, http.StatusInternalServerError, "")
		return
	}

	// Create a session cookie
	initial.GiveCookie(w, nickname)

	// Redirect to homepage
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func CheckName(name string) string {
	if name == "" {
		return "Name can't be null"
	}
	if len(name) > 20 {
		return "Name can't have more than 20 characters"
	}
	if len(name) < 4 {
		return "Name can't have less than 4 characters"
	}
	if strings.Contains(name, " ") {
		return "Name can't have a space in it"
	}

	return ""
}

// Hash a given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
