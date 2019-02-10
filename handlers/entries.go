package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func EntriesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Title: %v\n", vars["title"])
}
