package rest

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"gitlab.com/joshraphael/motdoftheday/pkg/apierror"
)

func (r Rest) DraftHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/draft.html"))
		vars := mux.Vars(req)
		method := apierror.MethodHTTP
		id := vars["post_id"]
		post_id, err := strconv.Atoi(id)
		if err != nil {
			msg := "invalid post_id in url: " + err.Error()
			log.Println(msg)
			apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", method)
			http.Error(w, msg, apiErr.Code())
			return
		}
		post, apiErr := r.processor.Draft(int64(post_id), apierror.MethodHTTP)
		if apiErr != nil {
			msg := "Error gathering draft posts: " + apiErr.Error()
			log.Println(msg)
			http.Error(w, msg, apiErr.Code())
			return
		}
		tmpl.Execute(w, post)
	}
}
