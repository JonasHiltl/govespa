package govespa

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/exp/maps"
)

type Query struct {
	client    *VespaClient
	iter      scanner
	ctx       context.Context
	yql       string
	options   QueryParameter
	variables url.Values
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

// AddVariables can be used to add any kind of key/value pair to the query.
// For example, AddVariable("ranking", "rank_albums") would select "rank_albums" as the Ranking Profile.
func (q *Query) AddVariable(key string, value string) *Query {
	q.variables.Add(key, value)
	return q
}

func (q *Query) AddParameter(p QueryParameter) *Query {
	q.options = p
	return q
}

// Get scans the first result into a destination.
// The destination needs to be a pointer to a struct which fields are annotated with the "vespa" Tag.
func (q *Query) Get(dest any) (QueryResponse, error) {
	res, err := q.fetch()
	if err != nil {
		return QueryResponse{}, err
	}
	if len(res.Root.Errors) > 0 {
		return QueryResponse{}, res.Root.Errors[0].ToError()
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

func (q *Query) fetch() (QueryResponse, error) {
	query := url.Values{}

	query.Add("yql", q.yql)
	maps.Copy(query, q.variables)
	maps.Copy(query, q.options.getQuery())

	resp, err := q.client.executeRequest(executeRequestParams{
		ctx:    q.ctx,
		path:   "/search/",
		query:  query,
		method: "GET",
		body:   nil,
	})
	if err != nil {
		return QueryResponse{}, err
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
