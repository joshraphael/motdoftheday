package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/joshraphael/diary/internal/server/rest"
	"gitlab.com/joshraphael/diary/pkg/database"
	"gitlab.com/joshraphael/diary/pkg/processors"
	"gitlab.com/joshraphael/diary/settings"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	v := validator.New()
	db_name := "./" + settings.DB_NAME
	if _, err := os.Stat(db_name); err != nil {
		msg := "Database " + db_name + " does not exist: " + err.Error()
		log.Fatal(msg)
	}
	db, err := sql.Open("sqlite3", db_name+"?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	database, err := database.New(db)
	if err != nil {
		log.Fatal(err)
	}
	processor := processors.New(database)
	apiHandler := rest.New(v, processor)
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
	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Serving at: " + addr)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
