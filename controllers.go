package main

import(
	"net/http"
	"fmt"
	"database/sql"
	"time"
	"strings"
	"os"
	"io"
	
    "github.com/julienschmidt/httprouter"
	
)

type NewImage struct {
	ImgURL string
	Caption string
	Tags string
	Imgname string
}

type Feed struct {
	Title string
	Content string
	Tags string
	Imgname string
	Postid string
	Userid string
}

type User struct {
	Userimg string
	Userid string
	Bio string
	LastActive string
	Name string
}

type ShowHideLinks struct {
	LoggedIn bool
	LoggedOut bool
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Println("path", r.URL.Path)
		fmt.Println(isLoggedIn(r))

	if r.Method == "GET" {
		if isLoggedIn(r) {
				http.Redirect(w, r, "/dash", http.StatusSeeOther)			
		} else {
		var lg ShowHideLinks;
		lg = ShowHideLinks{LoggedOut: true}
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", lg)
		tmpl.ExecuteTemplate(w, "index.html", nil)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
		}
	}
}


func Dash(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)
		
			row, error := db.Query("SELECT * from imgposts order by postid DESC LIMIT 8")
			if error != nil {
				panic(error.Error())
			}
		
		var title, content, tags, imgname, postid, userid string
		var res = []Feed{}

for row.Next() {
	switch err := row.Scan(&title, &content, &tags, &imgname, &postid, &userid); err {
		case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		case nil:
		//fmt.Println(title, content, tags, imgname, postid, userid)
		default:
		panic(err)
	}
	out := Feed{Title: title, Content: content, Tags: tags, Imgname: imgname, Postid: postid, Userid: userid}
		res = append(res, out)
		}
				var lg ShowHideLinks;
		lg = ShowHideLinks{LoggedIn: true}

		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", lg)
		tmpl.ExecuteTemplate(w, "dash.html", res)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
	
}  else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
type Tok struct {
	Token string
}
func NewPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if isLoggedIn(r) {
		t := time.Now()
        token := t.Format("20060102150405")
		tok := Tok{Token: token}
	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", nil)
		tmpl.ExecuteTemplate(w, "newimg.html", tok)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
	} else {
		fmt.Println("path", r.URL.Path)
		tt := r.FormValue("imgurl")
		ds := r.FormValue("cap")
		tg := r.FormValue("tags")
		token := r.FormValue("token")
		userid := "guest"
		//var res NewImage
		
       fmt.Println("method:", r.Method)
           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(err)
               return
           }
		   var s string = strings.ToLower(handler.Filename)
		   fileExtn  := s[len(s)-4:]
		   
		   if (fileExtn == ".jpg") || (fileExtn == ".png") || (fileExtn == ".gif") {		   
           defer file.Close()
           f, err := os.OpenFile("./img/"+token+fileExtn, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }
           defer f.Close()
           io.Copy(f, file)
       
		_, error := db.Exec("INSERT INTO imgposts(title, content, tags, imgname, postid, userid) VALUES(?,?,?,?,?,?)",tt, ds, tg, token+fileExtn, token, userid)
        if error != nil {
            panic(error.Error())
        }
		}
		http.Redirect(w, r, "/dash", http.StatusSeeOther)
	}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}


func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)
		var lg ShowHideLinks;
		lg = ShowHideLinks{LoggedOut: true}
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", lg)
		tmpl.ExecuteTemplate(w, "register.html", nil)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
	} else {
		fmt.Println("path", r.URL.Path)
		username := r.FormValue("uname")
		password := r.FormValue("pass")
		email := r.FormValue("email")
		uuid := generateNewUUID()
		
		if unameExists(username) {
		http.Redirect(w, r, "/register", 302)
		} else {
       fmt.Println("method:", r.Method)
		fmt.Println(uuid, username, password, email)
       
		_, error := db.Exec("INSERT INTO user(username, userid, password, email) VALUES(?,?,?,?)",username, uuid, password, email)
        if error != nil {
            panic(error.Error())
        }
		http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		}
	}



func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method == "GET" {
		fmt.Println("path", r.URL.Path)
		var lg ShowHideLinks;
		lg = ShowHideLinks{LoggedOut: true}
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", lg)
		tmpl.ExecuteTemplate(w, "login.html", nil)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
	} else {
		fmt.Println("path", r.URL.Path)
		username := r.FormValue("uname")
		password := r.FormValue("pass")
			fmt.Println("method:", r.Method)
			fmt.Println(username, password)
			var databaseUsername string
			var databasePassword string
			err := db.QueryRow("SELECT username, password FROM user WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
			fmt.Println("path", err)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
				
	        http.SetCookie(w, &http.Cookie{
                 Name:  "my-cookie", 
                 Value: username,  
                 Path: "/",
        })
	fmt.Println("Cookie set", username)
				
			http.Redirect(w, r, "/dash", http.StatusSeeOther)
		}
	}


func ShowPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		fmt.Println("path: ", r.URL.Path)
		
		postID := ps.ByName("postid")
			row, error := db.Query("SELECT * from imgposts where postid=?",postID)
			if error != nil {
				fmt.Println("I got nothing")
				panic(error.Error())
			}
		var title, content, tags, imgname, postid, userid string
		var res = []Feed{}

for row.Next() {
	switch err := row.Scan(&title, &content, &tags, &imgname, &postid, &userid); err {
		case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		case nil:
		//fmt.Println(title, content, tags, imgname, postid, userid)
		default:
		panic(err)
	}
	out := Feed{Title: title, Content: content, Tags: tags, Imgname: imgname, Postid: postid, Userid: userid}
		res = append(res, out)
		}
		
		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", nil)
		tmpl.ExecuteTemplate(w, "post.html", res)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
}

func generateNewUUID() string{
	 
		t := time.Now()
        newuserid := t.Format("20060102150405123")
		return newuserid
}


func UserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		fmt.Println("path: ", r.URL.Path)
		userID := ps.ByName("name")
			row, error := db.Query("SELECT * from users where userid=?",userID)
			if error != nil {
				fmt.Println("I got nothing")
				panic(error.Error())
			}
		var userimg, userid, bio, lastactive, name string
		var res = []User{}

for row.Next() {
	switch err := row.Scan(&userimg, &userid, &bio, &lastactive, &name); err {
		case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		case nil:
		//fmt.Println(userimg, userid, bio, lastactive, name)
		default:
		panic(err)
	}

	out := User{Userimg: userimg, Userid: userid, Bio: bio, LastActive: lastactive, Name: name}
		res = append(res, out)
		}

		tmpl.ExecuteTemplate(w, "head.html", nil)
		tmpl.ExecuteTemplate(w, "nav.html", nil)
		tmpl.ExecuteTemplate(w, "user.html", res)
		tmpl.ExecuteTemplate(w, "footer.html", nil)
}

func unameExists(uname string) bool{
var databaseUsername string
	err := db.QueryRow("SELECT username FROM user WHERE username=?", uname).Scan(&databaseUsername)
	fmt.Println("path", err)
	if err == nil {
		return true
	}
	return false
}


func DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		fmt.Println(":", r.Method)
		
		fmt.Println("path: ", r.URL.Path)
		postID := ps.ByName("postid")
				
			error := db.QueryRow("DELETE from imgposts where postid=?",postID)
			if error != nil {
				fmt.Println("I got nothing")
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
}

func isLoggedIn(r *http.Request)  bool {
	cookie, _ := r.Cookie("my-cookie")
    fmt.Println("Has: ", cookie)
	if cookie.Value == "" {
		return false
	} else {
		return true
	}
}


func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	        http.SetCookie(w, &http.Cookie{
                 Name:  "my-cookie", 
                 Value: "",  
                 Path: "/",
        })
	fmt.Println("Logged Out")
				
			http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
func SetCake(uname string) {
	        http.SetCookie(w, &http.Cookie{
                 Name:  "my-cookie", 
                 Value: uname,  
                 Path: "/",
        })
	fmt.Println("Cookie set", uname)
}

		// set cookie for storing token
		cookie := http.Cookie{}
		cookie.Name = "accessToken"
		cookie.Value = "ro8BS6Hiivgzy8Xuu09JDjlNLnSLldY5"
		cookie.MaxAge = time.Now().Minute() * 1
		cookie.Secure = true
		cookie.HttpOnly = true
		cookie.SameSite = http.SameSiteStrictMode
		http.SetCookie(w, &cookie)

func GetCake() string {
	c, err := req.Cookie("my-cookie") 
		if err != nil {  
				http.Error(w, http.StatusText(400), http.StatusBadRequest)  
		return 
		}  
	return c
}

*/

/*
func Create(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "create.html", nil)
}
*/
