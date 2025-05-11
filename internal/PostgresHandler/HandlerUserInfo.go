package PostgresHandler

import (
	"context"
	"database/sql"
	"fmt"
	"friend_graphql/graph/model"
	"friend_graphql/internal/logger"
	Postgres "friend_graphql/pkg/Posgres"
	"log"
	"strings"
)

func splitString(s string) []*string {
	postIDs := strings.Split(s[1:len(s)-1], ",")
	postIDSpointer := make([]*string, len(postIDs), len(postIDs))
	for i := 0; i < len(postIDs); i++ {
		postIDSpointer[i] = &postIDs[i]
	}
	return postIDSpointer
}

type Handler struct {
	DB *Postgres.Storage
}

type StorageHandler struct {
	storage *Postgres.Storage
}

func NewStorageHandler(storage *Postgres.Storage) *StorageHandler {
	return &StorageHandler{storage: storage}
}
func (s *StorageHandler) GetUserInfo(users []string) (map[string]model.UserInfo, error) {
	stmt, err := s.storage.Db.Prepare(getUserInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare GetUserInfo: %w", err)
	}
	ctx := context.Background()
	rows, err := stmt.QueryContext(ctx, users)
	if err != nil {
		return nil, fmt.Errorf("Poblem to return data GetUserInfo: %w", err)
	}
	usersInfo := make(map[string]model.UserInfo)
	var user model.UserInfo
	var userId string
	for rows.Next() {
		err = rows.Scan(&userId, &user.FirstName, &user.SecondName, &user.ImgUrl)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.GetLogger().Error("User does not exist")
				continue
			}
			return nil, fmt.Errorf("failed to scan rows GetUserInfo: %w", err)
		}
		usersInfo[userId] = user

	}
	return usersInfo, nil
}

func (s *StorageHandler) StorageGetUserInfoById(userID string, ctx context.Context) (*model.User, error) {
	user := new(model.User)
	var imagesStr string
	err := s.storage.Db.QueryRowContext(ctx, `SELECT user_id,first_name,second_name,img_url,images,birth_date,
       education,country,city FROM users_info where user_id=$1`, userID).Scan(&user.ID, &user.FirstName, &user.SecondName, &user.MainImgURL,
		&imagesStr, &user.BirthDate, &user.Education, &user.Country, &user.City)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	imagePointer := splitString(imagesStr)
	user.Images = imagePointer
	return user, nil
}
func (s *StorageHandler) StorageGetUserFriendIDs(user *model.User, ctx context.Context) error {
	rows, err := s.storage.Db.QueryContext(ctx, `SELECT friend_id,status FROM user_friends where user_id=$1`, user.ID)
	if err != nil {
		return err
	}
	subscribesIDS := []*string{}
	friendsIDs := []*string{}
	for rows.Next() {
		var friendID string
		var status bool
		err = rows.Scan(&friendID, &status)
		if err != nil {
			logger.GetLogger().Error("Failed to scan rows GetUserInfoById")
			continue
		}
		if status {
			friendsIDs = append(friendsIDs, &friendID)
		} else {
			subscribesIDS = append(subscribesIDS, &friendID)
		}
	}
	user.FriendIDs = friendsIDs
	user.SubscribesIDs = subscribesIDS
	return nil
}

func (s *StorageHandler) StorageGetUserFriendsAndSubscribers(userObj *model.User, ctx context.Context,
	limit, offset int32, friendStatus bool) ([]*model.User, error) {
	var friendIDs []*string
	var imagesStr string
	if friendStatus {
		friendIDs = userObj.FriendIDs
	} else {
		friendIDs = userObj.SubscribesIDs
	}
	result, err := s.storage.Db.QueryContext(ctx, `SELECT user_id,first_name,second_name,img_url,images,birth_date,
       education,country,city FROM users_info where user_id= ANY($1) LIMIT $2 OFFSET $3`, friendIDs, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0, len(friendIDs))
	for result.Next() {
		user := new(model.User)
		err = result.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.MainImgURL, &imagesStr,
			&user.BirthDate, &user.Education, &user.Country, &user.City)
		if err != nil {
			logger.GetLogger().Error("Failed to scan rows GetUserFriendsAndSubscribers")
			continue
		}
		imagePointer := splitString(imagesStr)
		user.Images = imagePointer
		users = append(users, user)
	}
	return users, nil
}
