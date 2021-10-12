package model

import "encoding/json"

// ElasticsearchGetRequestEnvelope is an envelope.
type ElasticsearchGetRequestEnvelope struct {
	Index        string          `json:"_index"`
	DocumentType string          `json:"_type"`
	DocumentID   string          `json:"_id"`
	Source       json.RawMessage `json:"_source"`
}

// ElasticsearchSearchRequestEnvelope is a model.
type ElasticsearchSearchRequestEnvelope struct {
	Took    int  `json:"took"`
	Timeout bool `json:"time_out"`
	Shards  struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source json.RawMessage `json:"_source"`
			Sort   []interface{}   `json:"sort,omitempty"`
		} `json:"hits"`
	} `json:"hits"`
}

// ElasticsearchBoolQuery is a model.
type ElasticsearchBoolQuery struct {
	From  *int `json:"from,omitempty"`
	Size  *int `json:"size,omitempty"`
	Query struct {
		Bool struct {
			Must               interface{} `json:"must,omitempty"`
			Should             interface{} `json:"should,omitempty"`
			Filter             interface{} `json:"filter,omitempty"`
			MustNot            interface{} `json:"must_not,omitempty"`
			MinimumShouldMatch interface{} `json:"minimum_should_match,omitempty"`
		} `json:"bool"`
	} `json:"query"`
	Sort interface{} `json:"sort,omitempty"`
}
