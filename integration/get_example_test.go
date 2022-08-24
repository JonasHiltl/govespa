package integration

import (
	"testing"

	"github.com/jonashiltl/govespa"
)

func TestGet(t *testing.T) {
	client := createClient()
	if client == nil {
		t.Fatal("Error creating Http Client")
	}

	c := govespa.NewClient(govespa.NewClientParams{
		HttpClient: client,
		BaseUrl:    "https://localhost:8090",
	})

	u := new(testUser)

	_, err := c.
		Get(govespa.DocumentId{
			Namespace:    "default",
			DocType:      "user",
			UserSpecific: "awgaw-1w234a-dw14ag-w1414a",
		}).
		Exec(u)
	if err != nil {
		t.Error(err)
	}

	if u.Firstname != "John" || u.Lastname != "Doe" {
		t.Errorf("Expected the user John Doe but got: %+v", u)
	}
}
