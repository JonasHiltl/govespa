package govespa

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type Remove struct {
	client *VespaClient
	ctx    context.Context
	id     DocumentId
	params RemoveParams
}

type RemoveParams struct {
	condition  string
	timeout    time.Duration
	route      string
	tracelevel uint8
}

func (r *Remove) WithContext(c context.Context) *Remove {
	r.ctx = c
	return r
}

func (r *Remove) AddParameter(p RemoveParams) *Remove {
	r.params = p
	return r
}

func (r *Remove) Exec() *vespaError {
	resp, err := r.client.executeRequest(executeRequestParams{
		ctx:    r.ctx,
		path:   r.id.toPath(),
		query:  r.params.getQuery(),
		method: "DELETE",
		body:   nil,
	})
	if err != nil {
		return fromError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		err := parseError(resp)
		return err
	}

	return nil
}

// TODO: implement "Delete Where" should be done through iteration

func (params RemoveParams) getQuery() (q url.Values) {
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
