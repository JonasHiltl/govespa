package integration

import (
	"reflect"
	"testing"

	"github.com/clubo-app/govespa"
)

type testUser struct {
	Username  string `vespa:"username"`
	Firstname string `vespa:"firstname"`
	Lastname  string `vespa:"lastname"`
}

func TestQuery(t *testing.T) {
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
		Query().
		AddYQL(`select * from user where default contains "joh"`).
		Get(u)
	if err != nil {
		t.Error(err)
	}

	exp := testUser{Username: "john.doe", Firstname: "John", Lastname: "Doe"}
	if !reflect.DeepEqual(*u, exp) {
		t.Errorf("Expected %+v to equal %+v", *u, exp)
	}
}
