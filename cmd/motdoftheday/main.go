package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/joshraphael/motdoftheday/internal/server/rest"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
	"gitlab.com/joshraphael/motdoftheday/pkg/processors"
	"gitlab.com/joshraphael/motdoftheday/settings"
	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	v := validator.New()
	db_name := "./" + settings.DB_NAME
	if _, err := os.Stat(db_name); err != nil {
		msg := "Database " + db_name + " does not exist: " + err.Error()
		log.Fatalln(msg)
	}
	db, err := sqlx.Open("sqlite3", db_name+"?_foreign_keys=on")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	d, err := database.New(db)
	if err != nil {
		log.Fatalln(err)
	}
	processor := processors.New(d)
	apiHandler := rest.New(v, processor)
	r.HandleFunc("/", apiHandler.HomeHandler).Methods("GET")
	r.HandleFunc("/drafts", apiHandler.DraftsHandler).Methods("GET")
	r.HandleFunc("/edit", apiHandler.EditHandler).Methods("GET")
	// Serve static files
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static").Handler(s).Methods("GET")

	// API
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/submit", apiHandler.SubmitHandler).Methods("POST")
	api.HandleFunc("/save", apiHandler.SaveHandler).Methods("POST")
	http.Handle("/", r)

	// Start HTTP Server
	addr := settings.HOST + ":" + settings.PORT
	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Serving at: " + addr)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
	log.Println("Server Shutdown Properly")
}
