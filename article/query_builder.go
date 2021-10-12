package article

import "fmt"

// BoolQuery is a concrete struct.
type BoolQuery struct {
	From  *int64 `json:"from,omitempty"`
	Size  *int64 `json:"size,omitempty"`
	Query struct {
		Bool struct {
			Must               []interface{} `json:"must,omitempty"`
			Should             []interface{} `json:"should,omitempty"`
			Filter             []interface{} `json:"filter,omitempty"`
			MustNot            []interface{} `json:"must_not,omitempty"`
			MinimumShouldMatch interface{}   `json:"minimum_should_match,omitempty"`
		} `json:"bool"`
	} `json:"query"`
	Sort []interface{} `json:"sort,omitempty"`
}

// NewBoolQuery is a constrctor.
func NewBoolQuery() *BoolQuery {
	return &BoolQuery{}
}

func (q *BoolQuery) AddFrom(from int64) *BoolQuery {
	q.From = &from
	return q
}

func (q *BoolQuery) AddSize(size int64) *BoolQuery {
	q.Size = &size
	return q
}

func (q *BoolQuery) AddFilterByAuthor(status string) *BoolQuery {
	filterByStatus := map[string]interface{}{
		"term": map[string]interface{}{
			"author": map[string]interface{}{
				"value": status,
			},
		},
	}

	q.Query.Bool.Must = append(q.Query.Bool.Must, filterByStatus)
	return q
}

func (q *BoolQuery) AddKeywordForSearch(keyword string) *BoolQuery {
	keyword = fmt.Sprintf(`"%s"`, keyword)
	search := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"type":  "bool_prefix",
			"query": keyword,
			"fields": []string{
				"title",
				"body",
				"voucherName",
			},
		},
	}

	q.Query.Bool.Must = append(q.Query.Bool.Must, search)

	return q
}

func (q *BoolQuery) AddSortByCreatedAt(directive string) *BoolQuery {
	sortByCreatedAt := map[string]interface{}{
		"created": map[string]interface{}{
			"order": directive,
		},
	}

	q.Sort = append(q.Sort, sortByCreatedAt)

	return q
}
