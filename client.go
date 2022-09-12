package govespa

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

type VespaClient struct {
	httpClient *http.Client
	headers    http.Header
	baseUrl    string
}

type NewClientParams struct {
	HttpClient *http.Client
	Headers    http.Header
	BaseUrl    string
}

func NewClient(params NewClientParams) *VespaClient {
	return &VespaClient{
		httpClient: params.HttpClient,
		headers:    params.Headers,
		baseUrl:    params.BaseUrl,
	}
}

func (v *VespaClient) Put(docId DocumentId) *Put {
	p := &Put{
		client: v,
		id:     docId,
		fields: make(map[string]any),
	}
	return p
}

func (v *VespaClient) Update(docId DocumentId) *Update {
	u := &Update{
		client: v,
		id:     docId,
		fields: make(map[string]map[string]any),
	}
	return u
}

func (v *VespaClient) Query() *Query {
	s := &Query{
		client:    v,
		variables: make(url.Values),
	}
	return s
}

func (v *VespaClient) Remove(docId DocumentId) *Remove {
	r := &Remove{
		client: v,
		id:     docId,
	}
	return r
}

func (v *VespaClient) Get(docId DocumentId) *Get {
	g := &Get{
		client: v,
		id:     docId,
	}
	return g
}

type executeRequestParams struct {
	ctx     context.Context
	path    string
	query   url.Values
	headers http.Header
	method  string
	body    io.Reader
}

func (v *VespaClient) executeRequest(params executeRequestParams) (*http.Response, error) {
	if params.ctx == nil {
		params.ctx = context.Background()
	}

	reqUrl, err := url.JoinPath(v.baseUrl, params.path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(params.ctx, params.method, reqUrl, params.body)
	if err != nil {
		return nil, err
	}

	dr, _ := httputil.DumpRequest(req, false)
	fmt.Println(string(dr))

	if params.query != nil {
		req.URL.RawQuery = params.query.Encode()
	}
	if params.headers != nil {
		req.Header = params.headers
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *VespaClient) getBaseUrl() (*url.URL, error) {
	u, err := url.Parse(v.baseUrl)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type OperationalParams struct {
	create     bool
	condition  string
	timeout    time.Duration
	route      string
	tracelevel uint8
}

func (params OperationalParams) getQuery() (q url.Values) {
	if params.create {
		q.Add("create", "true")
	}
	if params.condition != "" {
		q.Add("condition", params.condition)
	}
	if params.timeout != 0 {
		q.Add("timeout", strconv.FormatInt(params.timeout.Milliseconds(), 10))
	}
	if params.route != "" {
		q.Add("route", params.route)
	}
	if params.tracelevel != 0 {
		q.Add("tracelevel", strconv.FormatUint(uint64(params.tracelevel), 10))
	}
	return
}
