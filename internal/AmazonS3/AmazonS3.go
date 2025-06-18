package AmazonS3

import (
	"fmt"
	"friend_graphql/internal/config"
	"log"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type AmazonS3 struct {
	session      *session.Session
	bucketName   string
	domainServer string
}

func NewS3(cfg *config.Cfg) *AmazonS3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AmazonS3.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AmazonS3.AccessKey, cfg.AmazonS3.SecretKey, ""),
		Endpoint:    aws.String(cfg.AmazonS3.Endpoint),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &AmazonS3{session: sess, bucketName: cfg.AmazonS3.BucketName}
}

func (s *AmazonS3) UploadFile(file graphql.Upload) (string, error) {
	namePerm := strings.Split(file.Filename, ".")[1]
	img := file.File
	id := uuid.New().String()
	fmt.Println(namePerm)
	objectKey := "photos" + "/" + id + "." + namePerm
	_, err := s3manager.NewUploader(s.session).Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
		Body:   img,
	})
	if err != nil {
		return "", err
	}
	objectKey = s.domainServer + "/" + objectKey
	return objectKey, nil
}
