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
	Offset               uint64
	Hits                 uint32
	QueryProfile         string
	GroupingSessionCache *bool
	SearchChain          string
	Timeout              time.Duration
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

// Exec runs the Query but doesn't scan the fields of the result into a destination.
func (q *Query) Exec() (QueryResponse, error) {
	res, err := q.fetch()
	if err != nil {
		return QueryResponse{}, err
	}

	if len(res.Root.Errors) > 0 {
		return QueryResponse{}, res.Root.Errors[0].ToError()
	}
	return res, nil
}

// Get scans the fields of the first result/children into a destination.
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
	if p.Offset != 0 {
		q.Add("offset", strconv.FormatUint(p.Offset, 10))
	}
	if p.Hits != 0 {
		q.Add("hits", strconv.FormatUint(uint64(p.Hits), 10))
	}
	if p.QueryProfile != "" {
		q.Add("queryProfile", p.QueryProfile)
	}
	if p.GroupingSessionCache != nil {
		q.Add("groupingSessionCache", strconv.FormatBool(*p.GroupingSessionCache))
	}
	if p.SearchChain != "" {
		q.Add("searchChain", p.SearchChain)
	}
	if p.Timeout != 0 {
		q.Add("timeout", strconv.FormatInt(p.Timeout.Milliseconds(), 10))
	}
	return
}
