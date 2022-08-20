package main
import (
  	"fmt"
	"os"
	"html/template"
)

const PORT = "4000"
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	HandleRoutes()
	return
}


 func GetPort() string {
 	var port = os.Getenv("PORT")
 	if port == "" {
 		port = "4000"
 		fmt.Println("INFO: defaulting to " + port)
 	}
 	return ":" + port
}