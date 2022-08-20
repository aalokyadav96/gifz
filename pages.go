package main

import "net/http"

func Todo(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "todo.html", nil)
}
