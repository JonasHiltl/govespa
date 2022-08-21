package govespa

import (
	"bytes"
	"context"
	"encoding/json"
)

type Update struct {
	client *VespaClient
	ctx    context.Context
	id     DocumentId
	fields map[string]map[string]any
	params OperationalParams
}

func (u *Update) WithContext(c context.Context) *Update {
	u.ctx = c
	return u
}

func (u *Update) AddParameter(p OperationalParams) *Update {
	u.params = p
	return u
}

type updateBody struct {
	Fields    map[string]map[string]any `json:"fields"`
	Condition string                    `json:"condition"`
}

func (q *Update) Exec() error {
	b := updateBody{
		Fields:    q.fields,
		Condition: q.params.condition,
	}
	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	resp, err := q.client.executeRequest(executeRequestParams{
		ctx:    q.ctx,
		path:   q.id.toPath(),
		query:  q.params.getQuery(),
		method: "PUT",
		body:   bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res QueryResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	return nil
}

func (u *Update) Assign(field string, value any) *Update {
	u.safelyAddToMap("assign", field, value)
	return u
}

func (u *Update) Increment(field string, value any) *Update {
	u.safelyAddToMap("increment", field, value)
	return u
}

func (u *Update) Decrement(field string, value any) *Update {
	u.safelyAddToMap("decrement", field, value)
	return u
}

func (u *Update) Multiply(field string, value any) *Update {
	u.safelyAddToMap("multiply", field, value)
	return u
}

func (u *Update) Divide(field string, value any) *Update {
	u.safelyAddToMap("divide", field, value)
	return u
}

func (u *Update) Add(field string, value []any) *Update {
	u.safelyAddToMap("add", field, value)
	return u
}

// Adds a Weighted Set to a field. The weight can either be an Interger or a String.
// See documentation https://docs.vespa.ai/en/reference/document-json-format.html#add
func (u *Update) AddWeightedSet(field string, ws map[any]int) *Update {
	u.safelyAddToMap("add", field, ws)
	return u
}

// Remove can be used to remove an element from a map.
// See  https://docs.vespa.ai/en/reference/document-json-format.html#composite-remove
func (u *Update) Remove(field string) *Update {
	u.safelyAddToMap("remove", field, 0)
	return u
}

func (u *Update) RemoveWeightedSet(field string, element string) *Update {
	u.safelyAddToMap("remove", field, map[string]int{element: 0})
	return u
}

// Can be used to do an operation on a specific element in a weighted set or array.
// See https://docs.vespa.ai/en/reference/document-json-format.html#match
func (u *Update) Match(field string, element string, operation string, value any) *Update {
	u.fields[field]["match"] = map[string]any{
		"element": element,
		operation: value,
	}
	return u
}

// TODO: implement Add/Modify/Remove Method for tensors

func (u *Update) safelyAddToMap(operation string, field string, value any) {
	if u.fields[field] == nil {
		u.fields[field] = make(map[string]any)
	}
	u.fields[field][operation] = value
}
