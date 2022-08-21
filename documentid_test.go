package govespa

import "testing"

func TestDocIdToPath(t *testing.T) {
	tables := []struct {
		id   DocumentId
		path string
	}{
		{
			DocumentId{
				Namespace:    "first_namespace",
				DocType:      "my_doc_type",
				UserSpecific: "shakespeare",
			},
			"document/v1/first_namespace/my_doc_type/docid/shakespeare",
		},
		{
			DocumentId{
				Namespace:    "second_namespace",
				DocType:      "my_doc_type2",
				UserSpecific: "william",
			},
			"document/v1/second_namespace/my_doc_type2/docid/william",
		},
		{
			DocumentId{
				Namespace: "second_namespace",
				DocType:   "my_doc_type2",
			},
			"document/v1/second_namespace/my_doc_type2/docid/",
		},
	}

	for _, table := range tables {
		path := table.id.toPath()
		if path != table.path {
			t.Errorf("Path of docId %+v was incorrect, got: %v, expected: %v", table.id, path, table.path)
		}
	}
}

func TestDocIdToString(t *testing.T) {
	tables := []struct {
		id  DocumentId
		str string
	}{
		{
			DocumentId{
				Namespace:    "first_namespace",
				DocType:      "my_doc_type",
				UserSpecific: "shakespeare",
			},
			"id:first_namespace:my_doc_type::shakespeare",
		},
		{
			DocumentId{
				Namespace:    "second_namespace",
				DocType:      "my_doc_type2",
				UserSpecific: "william",
			},
			"id:second_namespace:my_doc_type2::william",
		},
		{
			DocumentId{
				Namespace: "second_namespace",
				DocType:   "my_doc_type2",
			},
			"id:second_namespace:my_doc_type2::",
		},
	}

	for _, table := range tables {
		str := table.id.String()
		if str != table.str {
			t.Errorf("String of %+v was incorrect, got: %v, expected: %v", table.id, str, table.str)
		}
	}

}

func TestParseDocId(t *testing.T) {
	tables := []struct {
		s  string
		id DocumentId
	}{
		{
			"id:first_namespace:my_doc_type::shakespeare",
			DocumentId{
				Namespace:    "first_namespace",
				DocType:      "my_doc_type",
				UserSpecific: "shakespeare",
			},
		},
		{
			"id:1:2:test:",
			DocumentId{
				Namespace: "1",
				DocType:   "2",
			},
		},
		{
			"id:nspace:type:skip:test",
			DocumentId{
				Namespace:    "nspace",
				DocType:      "type",
				UserSpecific: "test",
			},
		},
	}

	for _, table := range tables {
		id, _ := ParseDocId(table.s)
		if id != table.id {
			t.Errorf("Parsing \"%v\" failed, got: %+v, expected: %+v", table.s, id, table.id)
		}
	}
}

func BenchmarkParseDocId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseDocId("id:first_namespace:my_doc_type::shakespeare")
	}
}
