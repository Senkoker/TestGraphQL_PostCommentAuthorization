package model

type User struct {
	ID            string    `json:"id"`
	FirstName     string    `json:"firstName"`
	SecondName    string    `json:"secondName"`
	MainImgURL    string    `json:"mainImgUrl"`
	Images        []*string `json:"images,omitempty"`
	BirthDate     *string   `json:"birthDate,omitempty"`
	Education     *string   `json:"education,omitempty"`
	Country       *string   `json:"country,omitempty"`
	City          *string   `json:"city,omitempty"`
	FriendIDs     []*string `json:"friendIDs,omitempty"`
	SubscribesIDs []*string `json:"subscribesIDs,omitempty"`
	Friends       []*User   `json:"friends,omitempty"`
	Subscribes    []*User   `json:"subscribes,omitempty"`
	Posts         []*Post   `json:"posts,omitempty"`
}

type UserInfo struct {
	FirstName  string
	SecondName string
	ImgUrl     string
}

type UserFind struct {
	Find UserFindResult `json:"find"`
}

type UserFindOk struct {
	User *User `json:"user"`
}

func (UserFindOk) IsUserFindResult() {}

type UserID struct {
	Userid string `json:"userid"`
}
