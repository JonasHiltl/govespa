package govespa

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ErrorType int64

const (
	ConditionNotMet ErrorType = iota
	VespaFailure
	TransportFailure
)

type vespaError struct {
	Type       ErrorType `json:"type"`
	PathId     string    `json:"pathId"`
	Message    string    `json:"message"`
	VespaCode  int       `json:"code"`
	Source     string    `json:"source"`
	StatusCode int       `json:"status_code"`
}

func (e *vespaError) Error() string {
	return e.Message
}

func (e *vespaError) ToError() error {
	return errors.New(e.Message)
}

func fromError(e error) *vespaError {
	return &vespaError{
		Message: e.Error(),
	}
}

func parseError(resp *http.Response) *vespaError {
	e := vespaError{}
	switch resp.StatusCode {
	case 412:
		e.Type = ConditionNotMet
	case 502, 504, 507:
		e.Type = VespaFailure
	default:
		e.Type = TransportFailure
	}

	e.StatusCode = resp.StatusCode

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return fromError(err)
	}

	if err := json.Unmarshal(bb, &e); err != nil {
		return fromError(err)
	}

	return &e
}
