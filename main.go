package main

import (
	"diary/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HandleHome)
	http.ListenAndServe(":8080", nil)
}
