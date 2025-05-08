package domain

import (
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
)

type CommentDomain struct {
	StorageMongoDB MongoDBComementInterface
}

type MongoDBComementInterface interface {
	StorageGetComment(replyToID string, limit, offset int32) ([]*model.Comment, error)
}

func NewCommentDomain(storageMongoDB MongoDBComementInterface) *CommentDomain {
	return &CommentDomain{StorageMongoDB: storageMongoDB}
}

func (c *CommentDomain) GetComment(replyToID string, limit, offset int32) ([]*model.Comment, error) {
	op := "Get comment"
	commentLogger := logger.GetLogger().With("op", op)
	comments, err := c.StorageMongoDB.StorageGetComment(replyToID, limit, offset)
	if err != nil {
		commentLogger.Error("Get comment error", "err", err.Error())
		return nil, err
	}
	return comments, nil
}
