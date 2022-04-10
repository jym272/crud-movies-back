package main

import (
	"backend/models"
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq" // postgres driver
	"log"
	"net/http"
	"os"
	"time"
)

const (
	version = "0.0.1"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string //connection string
	}
}
type AppState struct {
	Status      string `json:"status"` //when rendering the json, this will be the key
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type Application struct {
	config config
	logger *log.Logger
	models *models.Models
}

func main() {
	config := config{}

	flag.IntVar(&config.port, "port", 8080, "Port to listen on")
	flag.StringVar(&config.env, "env", "development", "Application environment(development|production")
	//postgres://user:password@host/dbname?sslmode=disable
	flag.StringVar(&config.db.dsn, "db-dsn", "postgres://postgres:adini@localhost/movies?sslmode=disable", "Postgres Database connection string")
	flag.Parse()

	logger := log.New(os.Stdout, "Server: ", log.LstdFlags)
	db, err := openDB(config.db.dsn)
	if err != nil {
		logger.Fatal(err) //because we can't continue without a database
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Println(err)
		}
	}(db)

	app := &Application{
		config: config,
		logger: logger,
		models: models.NewModels(db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Listening on port %d", config.port)
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal(err) //because the server can't start
		return
	}

}

func openDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
