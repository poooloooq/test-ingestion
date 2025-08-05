package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/poooloooq/test-ingestion/cmd/config"
	"github.com/poooloooq/test-ingestion/cmd/repository"
	"github.com/poooloooq/test-ingestion/cmd/services"
)

func main() {

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	// Load environment variables from .env file or Secrets Manager
	config.Load(projectID)

	port := config.Config.Port
	http.HandleFunc("/posts", handleIngestion)
	http.HandleFunc("/posts/get", repository.GetPostsHandler)

	log.Printf("Server started on port:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func handleIngestion(w http.ResponseWriter, r *http.Request) {

	url := config.Config.APIURL
	source := config.Config.Source

	posts, err := services.GetAllPosts(url)
	if err != nil {
		log.Printf("Fetch Service error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	updated := services.ModifyPosts(posts, source)

	ctx := context.Background()
	if err := services.Insert(ctx, updated); err != nil {
		log.Printf("Storage Service error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Ingestion completed.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
