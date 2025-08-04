//go:build test

package services

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock types for Firestore
type MockFirestoreClient struct {
	mock.Mock
}

func (m *MockFirestoreClient) Collection(name string) *MockCollectionRef {
	args := m.Called(name)
	return args.Get(0).(*MockCollectionRef)
}
func (m *MockFirestoreClient) Close() error {
	return m.Called().Error(0)
}

type MockCollectionRef struct {
	mock.Mock
}

func (m *MockCollectionRef) Doc(id string) *MockDocumentRef {
	args := m.Called(id)
	return args.Get(0).(*MockDocumentRef)
}

type MockDocumentRef struct {
	mock.Mock
}

func (m *MockDocumentRef) Set(ctx context.Context, data interface{}) (interface{}, error) {
	args := m.Called(ctx, data)
	return args.Get(0), args.Error(1)
}

// Patchable Firestore client creator
var firestoreClientCreator = func(ctx context.Context, projectID string) (firestoreClient, error) {
	return firestore.NewClient(ctx, projectID)
}

// Interface for Firestore client to allow mocking
type firestoreClient interface {
	Collection(name string) *firestore.CollectionRef
	Close() error
}

// UpdatedPost mock struct
type UpdatedPost struct {
	ID   int
	Name string
}

// Insert function refactored for testability
func Insert(context context.Context, posts []UpdatedPost) error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := firestoreClientCreator(context, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	for _, post := range posts {
		_, err := client.Collection("posts").Doc(fmt.Sprint(post.ID)).Set(context, post)
		if err != nil {
			// log error, continue
		}
	}
	return nil
}

func TestInsert_Success(t *testing.T) {
	os.Setenv("GOOGLE_CLOUD_PROJECT", "test-project")
	mockClient := new(MockFirestoreClient)
	mockCol := new(MockCollectionRef)
	mockDoc := new(MockDocumentRef)

	firestoreClientCreator = func(ctx context.Context, projectID string) (firestoreClient, error) {
		return mockClient, nil
	}

	mockClient.On("Collection", "posts").Return(mockCol)
	mockCol.On("Doc", mock.Anything).Return(mockDoc)
	mockDoc.On("Set", mock.Anything, mock.Anything).Return(nil, nil)
	mockClient.On("Close").Return(nil)

	posts := []UpdatedPost{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
	err := Insert(context.Background(), posts)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockCol.AssertExpectations(t)
	mockDoc.AssertExpectations(t)
}

func TestInsert_ClientFailure(t *testing.T) {
	os.Setenv("GOOGLE_CLOUD_PROJECT", "test-project")
	firestoreClientCreator = func(ctx context.Context, projectID string) (firestoreClient, error) {
		return nil, errors.New("client error")
	}
	posts := []UpdatedPost{{ID: 1, Name: "A"}}
	err := Insert(context.Background(), posts)
	assert.Error(t, err)
}

func TestInsert_PostInsertFailure(t *testing.T) {
	os.Setenv("GOOGLE_CLOUD_PROJECT", "test-project")
	mockClient := new(MockFirestoreClient)
	mockCol := new(MockCollectionRef)
	mockDoc := new(MockDocumentRef)

	firestoreClientCreator = func(ctx context.Context, projectID string) (firestoreClient, error) {
		return mockClient, nil
	}

	mockClient.On("Collection", "posts").Return(mockCol)
	mockCol.On("Doc", mock.Anything).Return(mockDoc)
	mockDoc.On("Set", mock.Anything, mock.Anything).Return(nil, errors.New("insert error"))
	mockClient.On("Close").Return(nil)

	posts := []UpdatedPost{{ID: 1, Name: "A"}}
	err := Insert(context.Background(), posts)
	assert.NoError(t, err) // Should not fail, just log
}