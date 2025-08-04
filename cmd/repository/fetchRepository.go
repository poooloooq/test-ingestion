package repository

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/poooloooq/test-ingestion/cmd/config"
	"github.com/poooloooq/test-ingestion/cmd/services"
	"google.golang.org/api/iterator"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {

	//Load db config from env variables
	projectID := config.Config.GCProject
	ctx := context.Background()

	//Load firestore client
	//Needs creds.json for local runs
	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Printf("client Failure: %v", err)
		return
	}
	defer client.Close()

	iter := client.Collection("posts").Documents(ctx)

	var posts []services.UpdatedPost
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error reading from Firestore: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var post services.UpdatedPost
		if err := doc.DataTo(&post); err != nil {
			log.Printf("Failed to decode document: %v", err)
			http.Error(w, "Failed to decode data", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
