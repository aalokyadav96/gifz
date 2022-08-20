package main

import (
	"net/http"
  	"log"
	
    "github.com/julienschmidt/httprouter"
)

func HandleRoutes() {

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/dash", Dash)

	router.GET("/new", NewPhoto)
	router.POST("/new", NewPhoto)
	router.GET("/post/:postid", ShowPhoto)
	router.GET("/user/:name", UserProfile)
	router.GET("/delete/:postid", DeletePhoto)
	
	router.GET("/register", Register)
	router.POST("/register", Register)
	router.GET("/login", Login)
	router.POST("/login", Login)
	
	router.GET("/logout", Logout)
	
	router.NotFound = http.FileServer(http.Dir(""))
	router.ServeFiles("/usrimg/*filepath", http.Dir("usrimg"))
	router.ServeFiles("/img/*filepath", http.Dir("img"))
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.ServeFiles("/public/*filepath", http.Dir("public"))

err := http.ListenAndServe(GetPort(), router)
 	if err != nil {
		log.Fatal("error starting http server : ", router)
 	}

}

