package rest

import (
	"net/http"
	"text/template"
)

func (r Rest) HomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/home.html"))
		tmpl.Execute(w, nil)
	}
}
