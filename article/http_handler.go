package article

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jeri06/article-query-service/model"
	"github.com/jeri06/article-query-service/response"
	"github.com/sirupsen/logrus"
)

type HTTPHandler struct {
	Logger   *logrus.Logger
	Usecase  Usecase
	Validate *validator.Validate
}

func NewArticleHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, usecase Usecase) {
	handle := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/query-service/v1/article", handle.GetMany).Methods(http.MethodGet)
}

func (handler HTTPHandler) GetMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queryString := r.URL.Query()

	params := model.GetManyArticleParams{}
	params.Page, _ = strconv.ParseInt(queryString.Get("page"), 10, 64)
	params.Size, _ = strconv.ParseInt(queryString.Get("size"), 10, 64)
	params.Author = queryString.Get("author")
	params.Keyword = queryString.Get("keyword")

	if err := handler.validateRequestBody(params); err != nil {
		resp := response.NewErrorResponse(err, http.StatusBadRequest, nil, response.StatusInvalidPayload, err.Error())
		response.JSON(w, resp)
		return
	}

	resp := handler.Usecase.GetManyArticle(ctx, params)
	response.JSON(w, resp)
}

func (handler *HTTPHandler) validateRequestBody(body interface{}) (err error) {
	err = handler.Validate.Struct(body)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
