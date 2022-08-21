package integration

import (
	"context"
	"testing"

	"github.com/clubo-app/govespa"
)

func TestRemove(t *testing.T) {
	client := createClient()
	if client == nil {
		t.Fatal("Error creating Http Client")
	}

	c := govespa.NewClient(govespa.NewClientParams{
		HttpClient: client,
		BaseUrl:    "https://localhost:8090",
	})

	id := govespa.DocumentId{
		Namespace:    "default",
		DocType:      "user",
		UserSpecific: "awgaw-1w234a-dw14ag-w1414a",
	}

	err := c.
		Remove(id).
		WithContext(context.Background()).
		ByDocumentId()
	if err != nil {
		t.Error(err)
	}
}
