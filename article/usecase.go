package article

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/jeri06/article-query-service/entity"
	"github.com/jeri06/article-query-service/exception"
	"github.com/jeri06/article-query-service/model"
	"github.com/jeri06/article-query-service/response"
	"github.com/sirupsen/logrus"
)

const (
	RFC3339MillisWithTripleZeroFractionSecond string = "2006-01-02T15:04:05.000Z"
)

type Usecase interface {
	Save(ctx context.Context, payload model.Article) (resp response.Response)
	GetManyArticle(ctx context.Context, params model.GetManyArticleParams) (resp response.Response)
}

type usecase struct {
	serviceName string
	utcLoc      *time.Location
	logger      *logrus.Logger
	repository  Repository
}

func NewArticleUsecase(property UsecaseProperty) Usecase {
	return &usecase{
		serviceName: property.ServiceName,
		utcLoc:      property.UTCLoc,
		logger:      property.Logger,
		repository:  property.Repository,
	}
}

func (u usecase) Save(ctx context.Context, payload model.Article) (resp response.Response) {
	article := entity.Article{
		ID:      payload.ID,
		Author:  payload.Author,
		Title:   payload.Title,
		Body:    payload.Body,
		Created: payload.Created,
	}

	if err := u.repository.Save(ctx, article); err != nil {
		return response.NewErrorResponse(exception.ErrInternalServer, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	return response.NewSuccessResponse(nil, response.StatOK, "")

}

func (u usecase) GetManyArticle(ctx context.Context, params model.GetManyArticleParams) (resp response.Response) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	from := (params.Page - 1) * params.Size

	query := NewBoolQuery().
		AddFrom(from).
		AddSize(params.Size)

	if params.Author != "" {
		query.AddFilterByAuthor(params.Author)
	}

	if params.Keyword == "" {
		query.AddSortByCreatedAt("desc")
	} else {
		query.AddKeywordForSearch(params.Keyword)
	}

	bunchOfArticle, took, total, err := u.repository.FindMany(ctx, query)
	if err != nil {
		u.logger.Error(err.Error())
		if err == exception.ErrNotFound {
			return response.NewErrorResponse(err, http.StatusNotFound, nil, response.StatNotFound, "")
		}
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, "")
	}

	totalDataOnPage := len(bunchOfArticle)
	totalPage := int(math.Ceil(float64(total) / float64(params.Size)))

	meta := Meta{
		Took:            took,
		Page:            int(params.Page),
		TotalData:       total,
		TotalDataOnPage: totalDataOnPage,
		TotalPage:       totalPage,
	}

	return response.NewSuccessResponseWithMeta(bunchOfArticle, meta, response.StatOK, "")

}
