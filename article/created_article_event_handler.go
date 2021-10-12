package article

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/jeri06/article-query-service/model"
	"github.com/sirupsen/logrus"
)

type CreatedArticleEventHandler struct {
	logger   *logrus.Logger
	validate *validator.Validate
	usescase Usecase
}

func NewCreatedArticleEventHandler(logger *logrus.Logger, validate *validator.Validate, usecase Usecase) *CreatedArticleEventHandler {
	return &CreatedArticleEventHandler{logger, validate, usecase}
}

func (handler CreatedArticleEventHandler) Handle(ctx context.Context, message interface{}) (err error) {
	msg, ok := message.(*sarama.ConsumerMessage)
	if !ok {
		handler.logger.Error("Not a kafka message")
		return
	}

	var article model.Article

	if err = json.Unmarshal(msg.Value, &article); err != nil {
		handler.logger.Error(err)
		return
	}

	if err = handler.validateMessage(article); err != nil {
		handler.logger.Error(err)
		return
	}
	err = handler.usescase.Save(ctx, article).Error()

	return
}
func (handler CreatedArticleEventHandler) validateMessage(message interface{}) (err error) {
	err = handler.validate.Struct(message)
	if err == nil {
		return
	}

	errorFields := err.(validator.ValidationErrors)
	errorField := errorFields[0]
	err = fmt.Errorf("Invalid '%s' with value '%v'", errorField.Field(), errorField.Value())

	return
}
