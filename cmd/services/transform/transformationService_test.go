package transform

import (
	"testing"
	"time"

	"github.com/poooloooq/test-ingestion/cmd/services/fetch"
)

func TestModifyPosts(t *testing.T) {
	// Sample input posts
	posts := []fetch.Post{
		{
			UserID: 101,
			ID:     1,
			Title:  "Go Testing",
			Body:   "This is a test post",
		},
		{
			UserID: 102,
			ID:     2,
			Title:  "Another Post",
			Body:   "This is another post",
		},
	}

	source := "unit-test"

	// Call the function
	result := ModifyPosts(posts, source)

	// Check length
	if len(result) != len(posts) {
		t.Fatalf("Expected %d updated posts, got %d", len(posts), len(result))
	}

	// Check each field
	for i, updated := range result {
		expected := posts[i]

		if updated.UserID != expected.UserID {
			t.Errorf("UserID mismatch: expected %d, got %d", expected.UserID, updated.UserID)
		}
		if updated.ID != expected.ID {
			t.Errorf("ID mismatch: expected %d, got %d", expected.ID, updated.ID)
		}
		if updated.Title != expected.Title {
			t.Errorf("Title mismatch: expected %q, got %q", expected.Title, updated.Title)
		}
		if updated.Body != expected.Body {
			t.Errorf("Body mismatch: expected %q, got %q", expected.Body, updated.Body)
		}
		if updated.Source != source {
			t.Errorf("Source mismatch: expected %q, got %q", source, updated.Source)
		}
		if time.Since(updated.IngestedAt) > 2*time.Second {
			t.Errorf("IngestedAt %v is greater than current time", updated.IngestedAt)
		}
	}
}
