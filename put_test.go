package govespa

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindStruct(t *testing.T) {
	c := NewClient(NewClientParams{})

	tables := []struct {
		s   any
		exp map[string]any
	}{
		{
			s: struct {
				Username  string `vespa:"username"`
				Firstname string `vespa:"firstname"`
				Lastname  string `vespa:"lastname"`
				Hidden    string `vespa:"-"`
				Ignored   string `vespa:"ignored"`
			}{
				Username:  "john.doe",
				Firstname: "John",
				Lastname:  "Doe",
				Hidden:    "hide",
			},
			exp: map[string]any{
				"username":  "john.doe",
				"firstname": "John",
				"lastname":  "Doe",
			},
		},
		{
			s: struct {
				Age       int     `vespa:"age_field"`
				Relevance float64 `vespa:"relevance"`
			}{
				Age:       25,
				Relevance: 0.145,
			},
			exp: map[string]any{
				"age_field": 25,
				"relevance": 0.145,
			},
		},
		{
			s: struct {
				Tags []string `vespa:"tags"`
			}{
				Tags: []string{"club", "houseparty", "bar"},
			},
			exp: map[string]any{
				"tags": []string{"club", "houseparty", "bar"},
			},
		},
		{
			s: struct {
				Location struct {
					Lat float64 `vespa:"lat"`
					Lon float64 `vespa:"lon"`
				} `vespa:"location"`
			}{
				Location: struct {
					Lat float64 `vespa:"lat"`
					Lon float64 `vespa:"lon"`
				}{
					Lat: 54.23423,
					Lon: 8.1212,
				},
			},
			exp: map[string]any{
				"location": map[string]any{
					"lat": 54.23423,
					"lon": 8.1212,
				},
			},
		},
	}

	for _, table := range tables {
		p := c.
			Put(DocumentId{
				Namespace: "default",
				DocType:   "user",
			}).
			WithContext(context.Background()).
			BindStruct(table.s)

		assert.Equal(t, table.exp, p.fields, "Should be equal")
	}
}
