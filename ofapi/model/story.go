package model

import "time"

type Story struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"userId"`
	IsWatched bool        `json:"isWatched"`
	IsReady   bool        `json:"isReady"`
	Media     []Media     `json:"media"`
	CreatedAt time.Time   `json:"createdAt"`
	Question  interface{} `json:"question"`
	CanLike   bool        `json:"canLike"`
	IsLiked   bool        `json:"isLiked"`
}
