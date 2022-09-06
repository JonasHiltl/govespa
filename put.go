package govespa

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"

	"github.com/jmoiron/sqlx/reflectx"
)

type Put struct {
	client *VespaClient
	ctx    context.Context
	id     DocumentId
	fields map[string]any
	params OperationalParams
	mapper *reflectx.Mapper
}

func (p *Put) WithContext(c context.Context) *Put {
	p.ctx = c
	return p
}

func (p *Put) AddParameter(param OperationalParams) *Put {
	p.params = param
	return p
}

// BindStruct adds all values of the struct with the tag `vespa:"field_name"` to the fields object of the Put Request.
// use `vespa:"-" to exclude the value in the fields object`.
// Empty fields are ignored.
func (p *Put) BindStruct(s any) *Put {
	v := reflect.ValueOf(s)
	for v = reflect.ValueOf(s); v.Kind() == reflect.Ptr; {
		v = v.Elem()
	}
	fm := p.mapper.FieldMap(v)
	for k, v := range fm {
		if !v.IsZero() {
			p.fields[k] = v.Interface()
		}
	}
	return p
}

func (p *Put) BindMap(s map[string]any) *Put {
	// the manual loop allows chaining multiple BindMaps without loosing pre-binded key/values
	for k, v := range s {
		p.fields[k] = v
	}
	return p
}

type putBody struct {
	Fields    map[string]any `json:"fields"`
	Condition string         `json:"condition"`
}

func (p *Put) Exec() error {
	b := putBody{
		Fields:    p.fields,
		Condition: p.params.condition,
	}

	body, err := json.Marshal(b)
	if err != nil {
		return err
	}

	resp, err := p.client.executeRequest(executeRequestParams{
		ctx:     p.ctx,
		path:    p.id.toPath(),
		query:   p.params.getQuery(),
		headers: p.client.headers,
		method:  "POST",
		body:    bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		err := parseError(resp)
		return err.ToError()
	}

	return nil
}
