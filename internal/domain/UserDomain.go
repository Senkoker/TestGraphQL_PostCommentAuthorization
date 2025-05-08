package domain

import (
	"context"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
)

type UserDomain struct {
	StoragePostgres PostgresGetUserInfoInterface
}

func NewUserDomain(postgres PostgresGetUserInfoInterface) *UserDomain {
	return &UserDomain{StoragePostgres: postgres}
}

type PostgresGetUserInfoInterface interface {
	StorageGetUserInfoById(userID string, ctx context.Context) (*model.User, error)
	StorageGetUserFriend(user *model.User, ctx context.Context) error
	StorageGetUserFriendsAndSubscribers(userObj *model.User, ctx context.Context, limit, offset int32, friendStatus bool) ([]*model.User, error)
}

func (u *UserDomain) GetUserInfo(userID string) (*model.User, error) {
	op := "GetAllUserInfo"
	ctx := context.Background()
	userLogger := logger.GetLogger().With("op", op)
	user, err := u.StoragePostgres.StorageGetUserInfoById(userID, ctx)
	if err != nil {
		userLogger.Error("Error getting user info", "error", err.Error())
		return nil, err
	}
	err = u.StoragePostgres.StorageGetUserFriend(user, ctx)
	if err != nil {
		userLogger.Error("Error getting user friend", "error", err.Error())
		return nil, err
	}
	return user, nil
}

func (u *UserDomain) GetUserFriendSubscriber(user *model.User, ctx context.Context,
	limit, offset int32, friendStatus bool) ([]*model.User, error) {
	op := "GetUserFriendSubscriber"
	friendSubscriberLogger := logger.GetLogger().With("op", op)
	friends, err := u.StoragePostgres.StorageGetUserFriendsAndSubscribers(user, ctx, limit, offset, friendStatus)
	//Todo: проработать запрос так как отправляется несколько запросов одиночно
	if err != nil {
		friendSubscriberLogger.Error("Error getting user friends/subscribers", "error", err.Error())
		return nil, err
	}
	for _, friend := range friends {
		err = u.StoragePostgres.StorageGetUserFriend(friend, ctx)
		if err != nil {
			friendSubscriberLogger.Error("Error getting user friend", "error", err.Error())
			continue
		}
	}
	return friends, nil
}
