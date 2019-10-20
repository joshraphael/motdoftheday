package rest

import (
	"net/http"
	"text/template"
)

func (r Rest) DraftHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/draft.html"))

		tmpl.Execute(w, nil)
	}
}
