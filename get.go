package govespa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

var ErrNoDoc = errors.New("no document found")

type Get struct {
	client *VespaClient
	iter   scanner
	ctx    context.Context
	params GetParameter
	id     DocumentId
}

type GetParameter struct {
	cluster    string
	fieldSet   string
	timeout    time.Duration
	tracelevel uint8
}

type GetResponse struct {
	PathId string         `json:"pathId"`
	Id     string         `json:"id"`
	Fields map[string]any `json:"fields"`
}

func (g *Get) WithContext(c context.Context) *Get {
	g.ctx = c
	return g
}

func (r *Get) AddParameter(p GetParameter) *Get {
	r.params = p
	return r
}

func (g *Get) Exec(dest any) (GetResponse, error) {
	res, err := g.fetch()
	if err != nil {
		return GetResponse{}, err
	}
	if len(res.Fields) == 0 {
		return GetResponse{}, ErrNoDoc
	}

	i := scanner{
		res: []map[string]any{res.Fields},
	}
	err = i.Get(dest)
	if err != nil {
		return GetResponse{}, fromError(err)
	}

	return res, nil
}

func (g *Get) fetch() (GetResponse, error) {
	resp, err := g.client.executeRequest(executeRequestParams{
		ctx:    g.ctx,
		path:   g.id.toPath(),
		query:  g.params.getQuery(),
		method: "GET",
		body:   nil,
	})
	if err != nil {
		return GetResponse{}, err
	}
	defer resp.Body.Close()

	res := new(GetResponse)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return GetResponse{}, err
	}
	return *res, nil
}

func (p GetParameter) getQuery() (q url.Values) {
	if p.cluster != "" {
		q.Add("cluster", p.cluster)
	}
	if p.fieldSet != "" {
		q.Add("fieldSet", p.fieldSet)
	}
	if p.timeout != 0 {
		q.Add("timeout", fmt.Sprintf("%vs", p.timeout.Seconds()))
	}
	if p.tracelevel != 0 {
		q.Add("tracelevel", strconv.FormatUint(uint64(p.tracelevel), 10))
	}
	return
}
