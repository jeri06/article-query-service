package article

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/jeri06/article-query-service/entity"
	"github.com/jeri06/article-query-service/exception"
	"github.com/jeri06/article-query-service/model"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	Save(ctx context.Context, article entity.Article) (err error)
	FindMany(ctx context.Context, query interface{}) (bunchOfTicket []entity.Article, took, total int, err error)
}

type elasticsearchRequest interface {
	Do(ctx context.Context, transport esapi.Transport) (resp *esapi.Response, err error)
}

type repository struct {
	logger  *logrus.Logger
	client  *elasticsearch.Client
	index   string
	docType string
}

func NewArticleRepository(logger *logrus.Logger, client *elasticsearch.Client) Repository {
	var index string = "article"
	var docType string = "_doc"

	return &repository{
		logger:  logger,
		client:  client,
		index:   index,
		docType: docType,
	}
}
func (r *repository) do(ctx context.Context, req elasticsearchRequest) (responseBodyBuff []byte, err error) {
	var resp *esapi.Response
	if resp, err = req.Do(ctx, r.client); err != nil {
		r.logger.Error(err)
		err = exception.ErrInternalServer
		return
	}

	defer resp.Body.Close()

	if resp.IsError() {
		r.logger.Error(resp.String())
		if resp.StatusCode != http.StatusNotFound {
			err = exception.ErrInternalServer
			return
		}
		err = exception.ErrNotFound
		return
	}

	responseBodyBuff, _ = ioutil.ReadAll(resp.Body)
	return
}

func (r repository) Save(ctx context.Context, article entity.Article) (err error) {
	articleBuff := new(bytes.Buffer)
	json.NewEncoder(articleBuff).Encode(&article)

	req := esapi.IndexRequest{
		Index:        r.index,
		DocumentType: r.docType,
		DocumentID:   strconv.FormatInt(article.ID, 10),
		Body:         articleBuff,
	}

	_, err = r.do(ctx, &req)

	return
}
func (r repository) FindMany(ctx context.Context, query interface{}) (bunchOfArticle []entity.Article, took, total int, err error) {
	var responseBodyBuff []byte
	var dataLength int

	queryBuff := new(bytes.Buffer)
	if err = json.NewEncoder(queryBuff).Encode(query); err != nil {
		r.logger.Error(err)
		err = exception.ErrInternalServer
		return
	}

	req := esapi.SearchRequest{
		Index: []string{r.index},
		Body:  queryBuff,
	}

	if responseBodyBuff, err = r.do(ctx, &req); err != nil {
		return
	}

	var envelope model.ElasticsearchSearchRequestEnvelope
	json.Unmarshal(responseBodyBuff, &envelope)

	if dataLength = len(envelope.Hits.Hits); dataLength < 1 {
		err = exception.ErrNotFound
		return
	}

	total = envelope.Hits.Total.Value

	bunchOfArticle = make([]entity.Article, dataLength)
	for i, hit := range envelope.Hits.Hits {
		var article entity.Article
		json.Unmarshal(hit.Source, &article)
		bunchOfArticle[i] = article
	}

	return
}
