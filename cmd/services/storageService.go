package services

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/poooloooq/test-ingestion/cmd/config"
)

func Insert(context context.Context, posts []UpdatedPost) error {

	//Load db config from env variables
	projectID := config.Config.GCProject

	//Load firestore client
	//Needs creds.json for local runs
	client, err := firestore.NewClient(context, projectID)
	if err != nil {
		return fmt.Errorf("client Failure: %w", err)
	}
	defer client.Close()

	//Loop through slice and insert each post
	for _, post := range posts {
		_, err := client.Collection("posts").Doc(fmt.Sprint(post.ID)).Set(context, post)
		if err != nil {
			log.Printf("error inserting post Id %d : %v", post.ID, err)
		}
	}

	return nil
}
