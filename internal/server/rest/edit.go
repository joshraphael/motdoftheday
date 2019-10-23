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

func (r Rest) EditHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("./templates/edit.html"))
		vars := mux.Vars(req)
		method := apierror.MethodHTTP
		id := vars["post_history_id"]
		post_id, err := strconv.Atoi(id)
		if err != nil {
			msg := "invalid post_id in url: " + err.Error()
			log.Println(msg)
			apiErr := apierror.New(errors.New(msg), "BAD_REQUEST", method)
			http.Error(w, msg, apiErr.Code())
			return
		}
		post_history, apiErr := r.processor.Edit(int64(post_id), apierror.MethodHTTP)
		if apiErr != nil {
			msg := "Error gathering draft posts: " + apiErr.Error()
			log.Println(msg)
			http.Error(w, msg, apiErr.Code())
			return
		}
		tmpl.Execute(w, post_history)
	}
}
