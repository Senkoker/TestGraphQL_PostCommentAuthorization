package app

import (
	"friend_graphql/internal/AmazonS3"
	MongoDB2 "friend_graphql/internal/MongoDB"
	"friend_graphql/internal/PostgresHandler"
	producer "friend_graphql/internal/ProducerKafka"
	"friend_graphql/internal/Redis/RedisHandler"
	"friend_graphql/internal/config"
	"friend_graphql/internal/domain"
	"friend_graphql/internal/logger"
	"friend_graphql/internal/server"
	"friend_graphql/pkg/DBRedis"
	"friend_graphql/pkg/MongoDB"
	Postgres "friend_graphql/pkg/Posgres"
	"os"
	"os/signal"
	"syscall"
)

func App() {
	cfg := config.NewConfig()
	logger.LoggerInit(cfg.Logger.Debug)
	logger.GetLogger().Info("Config Message", *cfg)

	postgreStorage := Postgres.NewStorage(cfg.Postgres.Url)
	postgresHandler := PostgresHandler.NewStorageHandler(postgreStorage)

	rediStorage := DBRedis.NewRedis(cfg)
	redisHandler := RedisHandler.NewRedisHandler(rediStorage)

	mongoDBStorage := MongoDB.NewClientMongo(cfg)
	mongoHandler := MongoDB2.NewPostCommentHandler(mongoDBStorage)

	producerKafka := producer.NewProducer(cfg)

	amazonS3 := AmazonS3.NewS3(cfg)

	postDomain := domain.NewPostDomain(redisHandler, postgresHandler, amazonS3, producerKafka, mongoHandler)
	commentDomain := domain.NewCommentDomain(mongoHandler)
	userDomain := domain.NewUserDomain(postgresHandler)
	server := server.NewServer()
	server.QraphQLHandle(postDomain, commentDomain, userDomain)
	server.Start()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	producerKafka.Close()
	server.Stop()

}
