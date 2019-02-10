package main

import (
	"diary/handlers"
	"diary/settings"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/entries/{title}", handlers.EntriesHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	addr := settings.HOST + ":" + settings.PORT
	fmt.Println("Serving at: " + addr)
	http.ListenAndServe(addr, r)
}
