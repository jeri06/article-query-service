package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jeri06/article-query-service/article"
	"github.com/jeri06/article-query-service/config"
	"github.com/jeri06/article-query-service/consumer"
	"github.com/jeri06/article-query-service/response"
	"github.com/jeri06/article-query-service/server"
	_ "github.com/joho/godotenv/autoload" //
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

var (
	cfg          *config.Config
	utcLocation  *time.Location
	indexMessage string = "Application is running properly"
)

func init() {
	utcLocation, _ = time.LoadLocation("UTC")
	cfg = config.Load()
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)

	vld := validator.New()

	//kafka
	topics := "created-article-topic"
	ctx := context.Background()

	esClient, err := elasticsearch.NewClient(cfg.Elasticsearch)
	if err != nil {
		fmt.Println("elastic")
		logger.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/query-service", index)

	articleRepository := article.NewArticleRepository(logger, esClient)
	articleUsecase := article.NewArticleUsecase(article.UsecaseProperty{
		ServiceName: cfg.Application.Name,
		UTCLoc:      utcLocation,
		Logger:      logger,
		Repository:  articleRepository,
	})

	article.NewArticleHandler(logger, vld, router, articleUsecase)

	event := article.NewCreatedArticleEventHandler(logger, vld, articleUsecase)

	handler := cors.New(cors.Options{
		AllowedOrigins:   cfg.Application.AllowedOrigins,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	consumerHandler := consumer.NewConsumerHandler(event)

	client, err := sarama.NewConsumerGroup(cfg.SaramaKafka.Addresses, cfg.Application.Name, cfg.SaramaKafka.Config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumerHandler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

		}
	}()

	srv := server.NewServer(logger, handler, cfg.Application.Port)

	srv.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	<-sigterm

	srv.Close()

}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
