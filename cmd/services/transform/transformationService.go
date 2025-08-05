package transform

import (
	"time"

	"github.com/poooloooq/test-ingestion/cmd/services/fetch"
)

type UpdatedPost struct {
	fetch.Post
	IngestedAt time.Time `json:"ingested_at"`
	Source     string    `json:"source"`
}

func ModifyPosts(posts []fetch.Post, source string) []UpdatedPost {

	var UpdatedPosts []UpdatedPost

	//Cycle through slice of posts
	for _, p := range posts {

		//add new UpdatedPost struc to UpdatedPosts Slice for each Post struc
		UpdatedPosts = append(UpdatedPosts, UpdatedPost{
			Post:       p,
			IngestedAt: time.Now().UTC(),
			Source:     source,
		})
	}

	return UpdatedPosts
}
