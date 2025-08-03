package services

import (
	"time"

	"github.com/poooloooq/ingestion-pipeline/cmd/services/fetchService"
)

type EnrichedPost struct {
	fetchService.Post
	IngestedAt time.Time `json:"ingested_at"`
	Source     string    `json:"source"`
}

func EnrichPosts(posts []fetchService.Post, source string) []EnrichedPost {
	var enriched []EnrichedPost
	for _, p := range posts {
		enriched = append(enriched, EnrichedPost{
			Post:       p,
			IngestedAt: time.Now().UTC(),
			Source:     source,
		})
	}
	return enriched
}
