package main

import (
	"deploy-server/config"
	"deploy-server/handlers"
	"deploy-server/models"
	"deploy-server/utils"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load config
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config config.Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Initialize DB
	db, err := models.InitDB("db/deploy.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set up logging
	logFile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Set up routes
	http.HandleFunc("/deploy", utils.BasicAuth(handlers.DeployHandler(db, config)))

	// Start server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
