package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func getAllPosts(url string) ([]Post, error) {

	//Set Timeout from config
	timeout, err := time.ParseDuration(os.Getenv("HTTP_TIMEOUT"))
	if err != nil || timeout == 0 {
		timeout = 10 * time.Second
	}

	//Client Call
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)

	//Error Handling
	if err != nil {
		return nil, fmt.Errorf("Error while fetching: %w", err)
	}
	defer resp.Body.Close()

	//Response with unusual status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected Status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	posts := &[]Post{}
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("RUnable to parse the response: %w", err)
	}
	return *posts, nil
}
