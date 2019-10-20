package rest

import (
	"log"
	"net/http"
	"text/template"

	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
)

func (r Rest) DraftsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/drafts.html"))
		posts, apiErr := r.processor.Drafts(apierror.MethodHTTP)
		if apiErr != nil {
			msg := "Error gathering draft posts: " + apiErr.Error()
			log.Println(msg)
			http.Error(w, msg, apiErr.Code())
			return
		}
		tmpl.Execute(w, posts)
	}
}
