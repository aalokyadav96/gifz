package main

import (
	"net/http"
	"time"
	"fmt"
)



func Register(w http.ResponseWriter, r *http.Request) {
		uuid := generateNewUUID()
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", nil)
		tmpl.ExecuteTemplate(w, "register.html", nil)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
	} else {
		username := r.FormValue("uname")
		password := r.FormValue("pass")
		email := r.FormValue("email")
		if unameExists(username) {
		http.Redirect(w, r, "/register", 302)
		} else {
		db.Exec("INSERT INTO user (uuid, username, password, email) VALUES(?,?,?,?)", uuid, username, password, email)
		fmt.Println(uuid, username, password, email)
		http.Redirect(w, r, "/login", 302)
		}
	}
}
/*
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	} else {
		username := r.FormValue("uname")
		password := r.FormValue("pass")
		email := r.FormValue("email")
		//To-Do:  SELECT from table if it exists, SQL returns true
		db.Exec("SELECT * FROM user where username=? and password=?", username, password)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		
		http.Redirect(w, r, "/", 302)
}
*/
func generateNewUUID() string{
	 
		t := time.Now()
        newuserid := t.Format("20060102150405123456")
		return newuserid
}