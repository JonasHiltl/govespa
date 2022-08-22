package govespa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testUser struct {
	Username  string `vespa:"username"`
	Firstname string `vespa:"firstname"`
	Lastname  string `vespa:"lastname"`
	Age       int    `vespa:"age"`
}

func TestGet(t *testing.T) {
	tables := []struct {
		input []map[string]any
		exp   testUser
	}{
		{
			input: []map[string]any{{
				"username":  "John.Doe",
				"firstname": "John",
				"lastname":  "Doe",
				"age":       40,
			}},
			exp: testUser{Username: "John.Doe", Firstname: "John", Lastname: "Doe", Age: 40},
		},
		{
			input: []map[string]any{{
				"age": 18,
			}},
			exp: testUser{Age: 18},
		},
	}

	for _, table := range tables {
		dest := new(testUser)
		i := iter{
			res: table.input,
		}
		err := i.Get(dest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, table.exp, *dest, "Should be equal")
	}
}

func TestSelect(t *testing.T) {
	tables := []struct {
		input []map[string]any
		exp   []testUser
	}{
		{
			input: []map[string]any{
				{
					"username":  "John.Doe",
					"firstname": "John",
					"lastname":  "Doe",
					"age":       40,
				},
				{
					"username":  "will.shakespeare",
					"firstname": "William",
					"lastname":  "Shakespeare",
					"age":       34,
				},
			},
			exp: []testUser{
				{
					Username:  "John.Doe",
					Firstname: "John", Lastname: "Doe", Age: 40},
				{Username: "will.shakespeare", Firstname: "William", Lastname: "Shakespeare", Age: 34},
			},
		},
	}

	for _, table := range tables {
		dest := make([]testUser, 2)
		i := iter{
			res: table.input,
		}
		err := i.Select(&dest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, table.exp, dest, "Should be equal")
	}
}
