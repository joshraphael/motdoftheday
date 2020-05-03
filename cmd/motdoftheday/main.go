package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/joshraphael/motdoftheday/internal/server/rest"
	"gitlab.com/joshraphael/motdoftheday/pkg/config"
	"gitlab.com/joshraphael/motdoftheday/pkg/database"
	"gitlab.com/joshraphael/motdoftheday/pkg/processors"
	"gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
)

type Config struct {
	MotdOfTheDay config.Config `yaml:"motdoftheday" validate:"required"`
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	v := validator.New()
	conf := os.Getenv("CONFIG_ENV")
	cfg, err := initConfig(v, conf)
	db, sqlxDB, err := database.New(cfg.MotdOfTheDay.Database)
	if err != nil {
		log.Fatalln(err)
	}
	defer sqlxDB.Close()
	processor := processors.New(cfg.MotdOfTheDay.Processors, db)
	apiHandler := rest.New(v, processor)
	r.HandleFunc("/", apiHandler.HomeHandler).Methods("GET")
	r.HandleFunc("/drafts", apiHandler.DraftsHandler).Methods("GET")
	r.HandleFunc("/drafts/{post_id}", apiHandler.DraftHandler).Methods("GET")
	r.HandleFunc("/edit/{post_history_id}", apiHandler.EditHandler).Methods("GET")
	// Serve static files
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static").Handler(s).Methods("GET")

	// API
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/submit", apiHandler.SubmitHandler).Methods("POST")
	api.HandleFunc("/save", apiHandler.SaveHandler).Methods("POST")
	http.Handle("/", r)

	// Start HTTP Server
	addr := cfg.MotdOfTheDay.Rest.Host + ":" + cfg.MotdOfTheDay.Rest.Port
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

func initConfig(v *validator.Validate, file string) (*Config, error) {
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		msg := "Error reading config file '" + file + "': " + err.Error()
		return nil, errors.New(msg)
	}
	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		msg := "Error loading config '" + file + "': " + err.Error()
		return nil, errors.New(msg)
	}
	err = v.Struct(cfg)
	if err != nil {
		msg := "Error validating config '" + file + "': " + err.Error()
		return nil, errors.New(msg)
	}
	return &cfg, nil
}
