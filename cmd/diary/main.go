package main

import (
	"fmt"
	"net/http"

	"gitlab.com/joshraphael/diary/internal/server/rest"
	"gitlab.com/joshraphael/diary/settings"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	t := http.FileServer(http.Dir("./templates/"))
	r.PathPrefix("/static").Handler(s)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/submit", rest.SubmitHandler)
	r.PathPrefix("/").Handler(t)
	http.Handle("/", r)
	addr := settings.HOST + ":" + settings.PORT
	fmt.Println("Serving at: " + addr)
	http.ListenAndServe(addr, nil)
}
