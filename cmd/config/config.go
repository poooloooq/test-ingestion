package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APIURL      string
	Source      string
	Port        string
	GCProject   string
	HTTPTimeout string
}

var Config *AppConfig

func Load(projectID string) error {

	ctx := context.Background()

	// Check if .env file is present then load it
	if _, err := os.Stat(".env"); err == nil {

		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: error loading .env file: %v", err)
		}

		Config = &AppConfig{
			APIURL:      os.Getenv("API_URL"),
			Source:      os.Getenv("SOURCE"),
			Port:        os.Getenv("PORT"),
			HTTPTimeout: os.Getenv("HTTP_TIMEOUT"),
			GCProject:   projectID,
		}
		return nil
	}

	// Otherwise load secrets from Secret Manager
	apiURL, err := AccessSecret(ctx, projectID, "API_URL")
	if err != nil {
		return err
	}
	source, err := AccessSecret(ctx, projectID, "SOURCE")
	if err != nil {
		return err
	}
	port, err := AccessSecret(ctx, projectID, "PORT")
	if err != nil {
		port = "8080" //default port
	}
	httptimeout, err := AccessSecret(ctx, projectID, "HTTP_TIMEOUT")
	if err != nil {
		httptimeout = "10s" // default timeout
	}

	Config = &AppConfig{
		APIURL:      apiURL,
		Source:      source,
		Port:        port,
		HTTPTimeout: httptimeout,
		GCProject:   projectID,
	}

	return nil
}
