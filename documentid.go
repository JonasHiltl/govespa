package govespa

import (
	"errors"
	"strings"
)

type DocumentId struct {
	Namespace    string
	DocType      string
	UserSpecific string
	// TODO: support key/value pairs of DocumentID
}

// Current benchmarks on m1 Macbook
// goos: darwin
// goarch: arm64
// pkg: github.com/clubo-app/govespa
// BenchmarkParseDocId-8   	31781857	        37.31 ns/op

// TODO: cleean this mess up
func ParseDocId(s string) (DocumentId, error) {
	idx := strings.Index(s[0:], ":")
	if idx == -1 {
		return DocumentId{}, errors.New("Invalid format")
	}
	if s[0:idx] != "id" {
		return DocumentId{}, errors.New("No \"id\" prefix found")
	}

	subStr := s[idx+1:]
	nspcIdx := strings.Index(subStr, ":")
	if nspcIdx == -1 {
		return DocumentId{}, errors.New("No namespace found")
	}
	res := DocumentId{
		Namespace: subStr[:nspcIdx],
	}

	subStr = subStr[nspcIdx+1:]
	typeIdx := strings.Index(subStr, ":")
	if typeIdx == -1 {
		return DocumentId{}, errors.New("No docType found")
	}
	res.DocType = subStr[:typeIdx]

	subStr = subStr[typeIdx+1:]
	kvIdx := strings.Index(subStr, ":")
	if kvIdx == -1 {
		return DocumentId{}, errors.New("No key/values found")
	}

	res.UserSpecific = subStr[kvIdx+1:]

	return res, nil
}

func (i DocumentId) String() string {
	return "id:" + i.Namespace + ":" + i.DocType + "::" + i.UserSpecific
}

func (i DocumentId) toPath() string {
	var b strings.Builder
	b.WriteString("document/v1/")
	b.WriteString(i.Namespace)
	b.WriteString("/")
	b.WriteString(i.DocType)
	b.WriteString("/docid/")
	b.WriteString(i.UserSpecific)

	return b.String()
}
