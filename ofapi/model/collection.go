package model

import "fmt"

type Collection struct {
	CanAddUsers         bool         `json:"canAddUsers"`
	CanDelete           bool         `json:"canDelete"`
	CanManageUsers      bool         `json:"canManageUsers"`
	CanPinnedToChat     bool         `json:"canPinnedToChat"`
	CanPinnedToFeed     bool         `json:"canPinnedToFeed"`
	CanUpdate           bool         `json:"canUpdate"`
	Direction           any          `json:"direction"`
	ID                  any          `json:"id"`
	IsPinnedToChat      bool         `json:"isPinnedToChat"`
	IsPinnedToFeed      bool         `json:"isPinnedToFeed"`
	Name                string       `json:"name"`
	Order               string       `json:"order"`
	PostsCount          int          `json:"postsCount"`
	Posts               []any        `json:"posts"`
	CustomOrderUsersIds []any        `json:"customOrderUsersIds"`
	SortList            []any        `json:"sortList"`
	Type                string       `json:"type"`
	Users               []UserIDView `json:"users"`
	UsersCount          int          `json:"usersCount"`
}

func (c *Collection) StrID() string {
	return fmt.Sprintf("%v", c.ID)
}

func (c *Collection) IsStrTypeID() bool {
	if c.ID == nil {
		return false
	}
	_, ok := c.ID.(string)
	return ok
}
