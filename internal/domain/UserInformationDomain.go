package domain

import (
	"context"
	"friend_graphql/graph/model"
)

type UserInformation struct {
}
type StoragePostgresUser interface {
	GetUserInfoById(userID string, ctx context.Context) (*model.User, error)
}

func NewUserInformationDomain() *UserInformation {

}
