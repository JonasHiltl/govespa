package integration

import (
	"context"
	"testing"

	"github.com/jonashiltl/govespa"
)

func TestUpdate(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Fatal(err)
	}

	c := govespa.NewClient(govespa.NewClientParams{
		HttpClient: client,
		BaseUrl:    "https://localhost:8090",
	})

	err = c.
		Update(govespa.DocumentId{
			Namespace:    "default",
			DocType:      "user",
			UserSpecific: "awgaw-1w234a-dw14ag-w1414",
		}).
		WithContext(context.Background()).
		Assign("username", "new-john.doe").
		Exec()
	if err != nil {
		t.Error(err)
	}
}
