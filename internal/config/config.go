package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Cfg struct {
	KafkaInfo KafkaInfo
	Logger    LoggerDebug
	AmazonS3  AmazonS3
	Redis     Redis
	Postgres  PostgresInfo
	Mongo     MongoInfo
}

type PostgresInfo struct {
	Url string `env:"PG_URL"`
}
type KafkaInfo struct {
	KafkaAddresses string `env:"KAFKA_ADDRESSES"`
	KafkaTopicName string `env:"KAFKA_TOPIC_NAME"`
	KafkaPartition int32  `env:"KAFKA_PARTITION"`
	KafkaFlush     int    `env:"KAFKA_FLUSH"`
}
type LoggerDebug struct {
	Debug bool `env:"LoggerDebug"`
}

type MongoInfo struct {
	URL            string `env:"MONGO_URL"`
	Database       string `env:"MONGO_DATABASE"`
	CollectionName string `env:"MONGO_COLLECTION_NAME"`
}

type Redis struct {
	Address  string `env:"REDIS_ADDRESS" env-default:""`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB"`
}
type AmazonS3 struct {
	AccessKey    string `env:"SELECTEL_ACCESS_KEY" env-default:""`
	SecretKey    string `env:"SELECTEL_SECRET_KEY" env-default:""`
	BucketName   string `env:"SELECTEL_BUCKET_NAME" env-default:""`
	Region       string `env:"SELECTEL_REGION" env-default:""`
	Endpoint     string `env:"SELECTEL_ENDPOINT" env-default:""`
	DomainServer string `env:"SELECTEL_DOMAIN_SERVER" env-default:""`
}

func NewConfig() *Cfg {
	cfg := new(Cfg)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
