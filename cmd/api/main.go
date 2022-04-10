package main

import (
	"flag"
	"fmt"
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
}
type AppState struct {
	Status      string `json:"status"` //when rendering the json, this will be the key
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type Application struct {
	config config
	logger *log.Logger
}

func main() {
	config := config{}

	flag.IntVar(&config.port, "port", 8080, "Port to listen on")
	flag.StringVar(&config.env, "env", "development", "Application environment(development|production")
	flag.Parse()

	logger := log.New(os.Stdout, "Server: ", log.LstdFlags)

	app := &Application{config: config, logger: logger}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Listening on port %d", config.port)
	err := server.ListenAndServe()
	if err != nil {
		logger.Println(err)
		return
	}

}
