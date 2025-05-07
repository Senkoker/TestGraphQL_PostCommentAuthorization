package app

import (
	producer "friend_graphql/internal/ProducerKafka"
	"friend_graphql/internal/config"
	"friend_graphql/internal/logger"
	"friend_graphql/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func App() {
	cfg := config.NewConfig()
	logger.LoggerInit(cfg.Logger.Debug)
	producerKafka := producer.NewProducer(cfg)
	server := server.NewServer()
	server.QraphQLHandle(producerKafka)
	server.Start()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	producerKafka.Close()
	server.Stop()

}
