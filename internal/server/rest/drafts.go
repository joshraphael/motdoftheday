package rest

import (
	"net/http"
	"text/template"
)

func (r Rest) DraftsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/drafts.html"))
		tmpl.Execute(w, nil)
	}
}
