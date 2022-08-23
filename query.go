package govespa

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

type Query struct {
	client  *VespaClient
	iter    scanner
	ctx     context.Context
	yql     string
	options QueryParameter
}

// groupingSessionCache is a pointer so that we can distinguish between true/false/not defined
type QueryParameter struct {
	offset               uint64
	hits                 uint32
	queryProfile         string
	groupingSessionCache *bool
	searchChain          string
	timeout              time.Duration
}

func (q *Query) WithContext(c context.Context) *Query {
	q.ctx = c
	return q
}

func (q *Query) AddYQL(yql string) *Query {
	q.yql = yql
	return q
}

func (q *Query) AddParameter(p QueryParameter) *Query {
	q.options = p
	return q
}

// Get scans the first result into a destination.
// The destination needs to be a pointer to a struct which fields are annotated with the "vespa" Tag.
func (q *Query) Get(dest any) (QueryResponse, []vespaError) {
	res, vErr := q.fetch()
	if vErr != nil {
		return QueryResponse{}, []vespaError{*vErr}
	}
	if len(res.Root.Errors) > 0 {
		return QueryResponse{}, res.Root.Errors
	}

	if dest != nil {
		fields := getFieldsFromChildren(res.Root.Children)
		i := scanner{
			res: fields,
		}
		i.Get(dest)
	}
	return res, nil
}

func (q *Query) fetch() (QueryResponse, *vespaError) {
	query := url.Values{}
	q.options.addQueryToParams(query)
	query.Add("yql", q.yql)

	resp, err := q.client.executeRequest(executeRequestParams{
		ctx:    q.ctx,
		path:   "/search/",
		query:  query,
		method: "GET",
		body:   nil,
	})
	if err != nil {
		return QueryResponse{}, fromError(err)
	}
	defer resp.Body.Close()

	res := new(QueryResponse)
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return QueryResponse{}, fromError(err)
	}
	return *res, nil
}

func (p QueryParameter) getQuery() (q url.Values) {
	if p.offset != 0 {
		q.Add("offset", strconv.FormatUint(p.offset, 10))
	}
	if p.hits != 0 {
		q.Add("hits", strconv.FormatUint(uint64(p.hits), 10))
	}
	if p.queryProfile != "" {
		q.Add("queryProfile", p.queryProfile)
	}
	if p.groupingSessionCache != nil {
		q.Add("groupingSessionCache", strconv.FormatBool(*p.groupingSessionCache))
	}
	if p.searchChain != "" {
		q.Add("searchChain", p.searchChain)
	}
	if p.timeout != 0 {
		q.Add("timeout", strconv.FormatInt(p.timeout.Milliseconds(), 10))
	}
	return
}

func (p QueryParameter) addQueryToParams(params url.Values) {
	for k, v := range p.getQuery() {
		params.Add(k, v[0])
	}
}
