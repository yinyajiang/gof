package models

import (
	"time"
)

type SinglePost struct {
	ResponseType        string        `json:"responseType"`
	ID                  int           `json:"id"`
	PostedAt            time.Time     `json:"postedAt"`
	PostedAtPrecise     string        `json:"postedAtPrecise"`
	ExpiredAt           interface{}   `json:"expiredAt"`
	Author              Author        `json:"author"`
	Text                string        `json:"text"`
	RawText             string        `json:"rawText"`
	LockedText          bool          `json:"lockedText"`
	IsFavorite          bool          `json:"isFavorite"`
	CanReport           bool          `json:"canReport"`
	CanDelete           bool          `json:"canDelete"`
	CanComment          bool          `json:"canComment"`
	CanEdit             bool          `json:"canEdit"`
	IsPinned            bool          `json:"isPinned"`
	FavoritesCount      int           `json:"favoritesCount"`
	MediaCount          int           `json:"mediaCount"`
	IsMediaReady        bool          `json:"isMediaReady"`
	Voting              interface{}   `json:"voting"`
	IsOpened            bool          `json:"isOpened"`
	CanToggleFavorite   bool          `json:"canToggleFavorite"`
	StreamID            string        `json:"streamId"`
	Price               string        `json:"price"`
	HasVoting           bool          `json:"hasVoting"`
	IsAddedToBookmarks  bool          `json:"isAddedToBookmarks"`
	IsArchived          bool          `json:"isArchived"`
	IsPrivateArchived   bool          `json:"isPrivateArchived"`
	IsDeleted           bool          `json:"isDeleted"`
	HasURL              bool          `json:"hasUrl"`
	IsCouplePeopleMedia bool          `json:"isCouplePeopleMedia"`
	CommentsCount       int           `json:"commentsCount"`
	MentionedUsers      []interface{} `json:"mentionedUsers"`
	LinkedUsers         []interface{} `json:"linkedUsers"`
	TipsAmount          string        `json:"tipsAmount"`
	TipsAmountRaw       string        `json:"tipsAmountRaw"`
	Media               []Medium      `json:"media"`
	CanViewMedia        bool          `json:"canViewMedia"`
	Preview             []interface{} `json:"preview"`
}

type Author struct {
	ID   int    `json:"id"`
	View string `json:"_view"`
}

type Files struct {
	Full          Full          `json:"full"`
	Thumb         Thumb         `json:"thumb"`
	Preview       Preview       `json:"preview"`
	SquarePreview SquarePreview `json:"squarePreview"`
	Drm           Drm           `json:"drm"`
}

type Full struct {
	URL     string        `json:"url"`
	Width   int           `json:"width"`
	Height  int           `json:"height"`
	Size    int           `json:"size"`
	Sources []interface{} `json:"sources"`
}

type SquarePreview struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
}

type Thumb struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
}

type Info struct {
	Source  Source  `json:"source"`
	Preview Preview `json:"preview"`
}

type Medium struct {
	ID               int64        `json:"id"`
	Type             string       `json:"type"`
	ConvertedToVideo bool         `json:"convertedToVideo"`
	CanView          bool         `json:"canView"`
	HasError         bool         `json:"hasError"`
	CreatedAt        *time.Time   `json:"createdAt"`
	Info             Info         `json:"info"`
	Source           Source       `json:"source"`
	SquarePreview    string       `json:"squarePreview"`
	Full             string       `json:"full"`
	Preview          string       `json:"preview"`
	Thumb            string       `json:"thumb"`
	HasCustomPreview bool         `json:"hasCustomPreview"`
	Files            Files        `json:"files"`
	VideoSources     VideoSources `json:"videoSources"`
}

type Preview struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int    `json:"size"`
	URL    string `json:"url"`
}

type Source struct {
	Source   string `json:"source"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Duration int    `json:"duration"`
}

type VideoSources struct {
	Resolution720 interface{} `json:"720"`
	Resolution240 interface{} `json:"240"`
}

type Dash struct {
	CloudFrontPolicy    string `json:"CloudFront-Policy"`
	CloudFrontSignature string `json:"CloudFront-Signature"`
	CloudFrontKeyPairID string `json:"CloudFront-Key-Pair-Id"`
}

type Drm struct {
	Manifest  Manifest  `json:"manifest"`
	Signature Signature `json:"signature"`
}

type Hls struct {
	CloudFrontPolicy    string `json:"CloudFront-Policy"`
	CloudFrontSignature string `json:"CloudFront-Signature"`
	CloudFrontKeyPairID string `json:"CloudFront-Key-Pair-Id"`
}

type Manifest struct {
	Hls  *string `json:"hls"`
	Dash *string `json:"dash"`
}

type Signature struct {
	Hls  Hls  `json:"hls"`
	Dash Dash `json:"dash"`
}
