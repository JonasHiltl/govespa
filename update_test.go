package govespa

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestAssign(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value any
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "username", value: "john.doe"},
			exp: map[string]map[string]any{
				"username": {
					"assign": "john.doe",
				},
			},
		},

		{
			in: input{field: "age", value: 18},
			exp: map[string]map[string]any{
				"age": {
					"assign": 18,
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Assign(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestIncrement(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value any
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "age", value: 0.145},
			exp: map[string]map[string]any{
				"age": {
					"increment": 0.145,
				},
			},
		},

		{
			in: input{field: "age", value: 18},
			exp: map[string]map[string]any{
				"age": {
					"increment": 18,
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Increment(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestDecrement(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value any
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "age", value: 0.145},
			exp: map[string]map[string]any{
				"age": {
					"decrement": 0.145,
				},
			},
		},

		{
			in: input{field: "age", value: 18},
			exp: map[string]map[string]any{
				"age": {
					"decrement": 18,
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Decrement(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestMultiply(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value any
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "age", value: 0.145},
			exp: map[string]map[string]any{
				"age": {
					"multiply": 0.145,
				},
			},
		},

		{
			in: input{field: "age", value: 18},
			exp: map[string]map[string]any{
				"age": {
					"multiply": 18,
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Multiply(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestDivide(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value any
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "age", value: 0.145},
			exp: map[string]map[string]any{
				"age": {
					"divide": 0.145,
				},
			},
		},

		{
			in: input{field: "age", value: 18},
			exp: map[string]map[string]any{
				"age": {
					"divide": 18,
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Divide(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestAdd(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value []any
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "tracks", value: []any{"Lay Lady Lay", "Every Grain of Sand"}},
			exp: map[string]map[string]any{
				"tracks": {
					"add": []any{"Lay Lady Lay", "Every Grain of Sand"},
				},
			},
		},

		{
			in: input{field: "age", value: []any{1, 2, 3, 4}},
			exp: map[string]map[string]any{
				"age": {
					"add": []any{1, 2, 3, 4},
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Add(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestAddWeightedSet(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field string
		value map[any]int
	}

	tables := []struct {
		in  input
		exp map[string]map[string]any
	}{
		{
			in: input{field: "int_weighted_set", value: map[any]int{123: 123, 456: 100}},
			exp: map[string]map[string]any{
				"int_weighted_set": {
					"add": map[any]int{123: 123, 456: 100},
				},
			},
		},

		{
			in: input{field: "string_weighted_set", value: map[any]int{"item 1": 144, "item 2": 7}},
			exp: map[string]map[string]any{
				"string_weighted_set": {
					"add": map[any]int{"item 1": 144, "item 2": 7},
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			AddWeightedSet(table.in.field, table.in.value)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestRemove(t *testing.T) {
	c := NewClient(NewClientParams{})

	tables := []struct {
		field string
		exp   map[string]map[string]any
	}{
		{
			field: "string_map{item 2}",
			exp: map[string]map[string]any{
				"string_map{item 2}": {
					"remove": 0,
				},
			},
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			Remove(table.field)

		if !reflect.DeepEqual(u.fields, table.exp) {
			t.Errorf("Expected %+v to equal %+v", u.fields, table.exp)
		}
	}
}

func TestRemoveWeightedSet(t *testing.T) {
	c := NewClient(NewClientParams{})

	type input struct {
		field   string
		element string
	}

	tables := []struct {
		in  input
		exp any
	}{
		{
			in:  input{field: "string_weighted_set", element: "item 2"},
			exp: "{\"string_weighted_set\":{\"remove\":{\"item 2\":0}}}",
		},
	}

	for _, table := range tables {
		u := c.
			Update(DocumentId{}).
			RemoveWeightedSet(table.in.field, table.in.element)

		json, err := json.Marshal(u.fields)
		if err != nil {
			t.Error(err)
		}
		if string(json) != table.exp {
			t.Errorf("Expected %v to equal %v", string(json), table.exp)
		}
	}
}
