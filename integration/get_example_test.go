package integration

import (
	"log"
	"testing"

	"github.com/jonashiltl/govespa"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	client, err := createClient()
	if err != nil {
		log.Println("Create Client Error")
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
	exp := testUser{Username: "john.doe", Firstname: "John", Lastname: "Doe"}

	// make sure the document we later want to "get" is inserted
	err = c.Put(id).BindStruct(exp).Exec()
	if err != nil {
		t.Fatal(err)
	}

	u := new(testUser)
	_, err = c.
		Get(id).
		Exec(u)
	if err != nil {
		log.Println("Get error")
		t.Fatal(err)
	}

	assert.Equal(t, exp, *u, "Should be equal")
}
