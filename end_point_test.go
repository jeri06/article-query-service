package main_test

import (
	"fmt"
	"io/ioutil"

	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jeri06/article-query-service/article"
	"github.com/jeri06/article-query-service/config"
	"github.com/jeri06/article-query-service/model"

	"github.com/sirupsen/logrus"
)

func Test_users(t *testing.T) {
	utcLocation, _ := time.LoadLocation("UTC")
	cfg := config.Load()
	r := mux.NewRouter()
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	vld := validator.New()

	esClient, err := elasticsearch.NewClient(cfg.Elasticsearch)
	if err != nil {
		fmt.Println("elastic")
		logger.Fatal(err)
	}

	articleRepository := article.NewArticleRepository(logger, esClient)
	articleUsecase := article.NewArticleUsecase(article.UsecaseProperty{
		ServiceName: cfg.Application.Name,
		UTCLoc:      utcLocation,
		Logger:      logger,
		Repository:  articleRepository,
	})

	//kafka
	hh := article.HTTPHandler{
		Logger:   logger,
		Validate: vld,
		Usecase:  articleUsecase,
	}
	params := model.GetManyArticleParams{
		Page:    1,
		Size:    10,
		Keyword: "",
		Author:  "",
	}

	r.HandleFunc("/query-service/v1/article", hh.GetMany)

	req := httptest.NewRequest(http.MethodGet, "/query-service/v1/article", nil)
	req.Header.Set("Content-Type", "application/json")
	query := req.URL.Query()
	query.Add("page", strconv.FormatInt(int64(params.Page), 10))
	query.Add("size", strconv.FormatInt(int64(params.Size), 10))
	query.Add("author", params.Author)
	query.Add("keyword", params.Keyword)

	req.URL.RawQuery = query.Encode()

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	res := rr.Result()
	_, _ = ioutil.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, rr.Code)

}
