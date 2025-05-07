package ProducerKafka

import (
	"errors"
	"friend_graphql/internal/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

type Producer struct {
	producer   *kafka.Producer
	partition  int32
	topic      *string
	kafkaFlush int
}

func NewProducer(cfg *config.Cfg) *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaInfo.KafkaAddresses,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Producer{producer: producer, partition: cfg.KafkaInfo.KafkaPartition, topic: &cfg.KafkaInfo.KafkaTopicName,
		kafkaFlush: cfg.KafkaInfo.KafkaFlush}
}

func (p *Producer) Produce(msg []byte) error {
	kafkaResponse := make(chan kafka.Event)
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     p.topic,
			Partition: p.partition,
		},
		Value:         msg,
		Key:           nil,
		Timestamp:     time.Time{},
		TimestampType: 0,
		Headers:       nil,
	}, kafkaResponse)
	event := <-kafkaResponse
	switch event.(type) {
	case kafka.Error:
		return err
	case *kafka.Message:
		return nil
	default:
		return errors.New("Unknown Kafka Message")
	}
}

func (p *Producer) Close() {
	p.producer.Flush(p.kafkaFlush)
	p.producer.Close()
}
