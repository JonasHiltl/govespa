package integration

import (
	"context"
	"testing"

	"github.com/jonashiltl/govespa"
)

func TestRemove(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatal(err)
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

	err = c.
		Remove(id).
		WithContext(context.Background()).
		Exec()
	if err != nil {
		t.Error(err)
	}
}
