package govespa

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/mitchellh/mapstructure"
)

type Put struct {
	client *VespaClient
	ctx    context.Context
	id     DocumentId
	fields map[string]any
	params OperationalParams
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
//
// Use `vespa:"-"` to exclude the value in the fields object.
// Use `vespa:",omitempty"` to ignore the field if values if the default value.
func (p *Put) BindStruct(s any) *Put {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &p.fields,
		TagName: "vespa",
	})
	if err != nil {
		log.Println(err)
	}

	err = d.Decode(s)
	if err != nil {
		log.Println(err)
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
