package govespa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

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

func (g *Get) Exec(dest any) (GetResponse, *vespaError) {
	resp, err := g.client.executeRequest(executeRequestParams{
		ctx:    g.ctx,
		path:   g.id.toPath(),
		query:  g.params.getQuery(),
		method: "GET",
		body:   nil,
	})
	if err != nil {
		return GetResponse{}, fromError(err)
	}
	defer resp.Body.Close()

	res := new(GetResponse)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return GetResponse{}, fromError(err)
	}

	i := scanner{
		res: []map[string]any{res.Fields},
	}
	err = i.Get(dest)
	if err != nil {
		return GetResponse{}, fromError(err)
	}
	return *res, nil
}

// TODO: implement visiting with continuation token through the iterator
func (g *Get) Visit() {

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
