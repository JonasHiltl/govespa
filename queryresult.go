package govespa

type QueryResponse struct {
	Root Root `json:"root"`
}

type Root struct {
	Id        string       `json:"id,omitempty"`
	Relevance float64      `json:"relevance"`
	Label     string       `json:"label,omitempty"`
	Source    string       `json:"source,omitempty"`
	Value     string       `json:"value,omitempty"`
	Types     []string     `json:"types,omitempty"`
	Children  []children   `json:"children"`
	Fields    fields       `json:"fields"`
	Coverage  coverage     `json:"coverage"`
	Limits    limits       `json:"limits"`
	Errors    []vespaError `json:"errors,omitempty"`
}

type children struct {
	Id        string         `json:"id,omitempty"`
	Relevance float64        `json:"relevance"`
	Label     string         `json:"label,omitempty"`
	Source    string         `json:"source,omitempty"`
	Fields    map[string]any `json:"fields"`
}

func getFieldsFromChildren(ch []children) []map[string]any {
	fields := make([]map[string]any, len(ch))
	for i, c := range ch {
		fields[i] = c.Fields
	}
	return fields
}

type coverage struct {
	Coverage    int8     `json:"coverage"`
	Documents   int64    `json:"documents"`
	Full        bool     `json:"bool"`
	Nodes       int      `json:"nodes"`
	Results     int      `json:"results"`
	ResultsFull int      `json:"resultsFull"`
	Degraded    degraded `json:"degraded"`
}

type degraded struct {
	MatchPhase      bool `json:"match-phase"`
	Timeout         bool `json:"timeout"`
	AdaptiveTimeout bool `json:"adaptive-timeout"`
	NonIdealState   bool `json:"non-ideal-state"`
}

type fields struct {
	Summaryfeatures string `json:"summaryfeatures,omitempty"`
	Matchfeatures   string `json:"matchfeatures,omitempty"`
	TotalCount      int    `json:"totalCount"`
}

type limits struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
