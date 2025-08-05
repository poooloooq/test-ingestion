package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/poooloooq/test-ingestion/cmd/config"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func GetAllPosts(url string) ([]Post, error) {

	//Set Timeout from config
	timeout, err := time.ParseDuration(config.Config.HTTPTimeout)
	if err != nil || timeout == 0 {
		timeout = 10 * time.Second
	}

	//Client Call
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)

	//Error Handling
	if err != nil {
		return nil, fmt.Errorf("error while fetching: %w", err)
	}
	defer resp.Body.Close()

	//Response with unusual status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	//Read from response and parse JSON to slice of posts
	body, _ := io.ReadAll(resp.Body)
	posts := &[]Post{}
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("unable to parse the response: %w", err)
	}
	return *posts, nil
}
