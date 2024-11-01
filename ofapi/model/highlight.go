package model

import "time"

type Highlight struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	Title        string    `json:"title"`
	CoverStoryID int64     `json:"coverStoryId"`
	Cover        string    `json:"cover"`
	StoriesCount int       `json:"storiesCount"`
	CreatedAt    time.Time `json:"createdAt"`
	Stories      []Story   `json:"stories"`
}
