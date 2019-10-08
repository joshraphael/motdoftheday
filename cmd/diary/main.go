package main

import (
	"fmt"
	"net/http"

	"gitlab.com/joshraphael/diary/internal/server/rest"
	"gitlab.com/joshraphael/diary/settings"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	v := validator.New()
	apiHandler := rest.New(v)
	// Serve static files
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static").Handler(s)

	// API
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/submit", apiHandler.SubmitHandler)
	api.HandleFunc("/save", apiHandler.SaveHandler)

	// Render templates
	t := http.FileServer(http.Dir("./templates/"))
	r.PathPrefix("/").Handler(t)
	http.Handle("/", r)

	// Start HTTP Server
	addr := settings.HOST + ":" + settings.PORT
	fmt.Println("Serving at: " + addr)
	http.ListenAndServe(addr, nil)
}
