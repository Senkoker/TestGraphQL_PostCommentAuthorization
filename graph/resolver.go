package graph

import (
	"friend_graphql/graph/model"
	"github.com/99designs/gqlgen/graphql"
)

type ProducerKafkaInterface interface {
	Produce(msg []byte) error
}

type PostgresUserInfo interface {
	GetUserInfo(users []string) (map[string]model.UserInfo, error)
}

type AmazonS3Interface interface {
	UploadFile(file graphql.Upload) (string, error)
}
type Resolver struct {
	producer ProducerKafkaInterface
	amazonS3 AmazonS3Interface
	postgres PostgresUserInfo
}
