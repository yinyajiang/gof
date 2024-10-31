package model

import "time"

type Post struct {
	ResponseType        string     `json:"responseType"`
	ID                  int        `json:"id"`
	PostedAt            time.Time  `json:"postedAt"`
	PostedAtPrecise     string     `json:"postedAtPrecise"`
	ExpiredAt           any        `json:"expiredAt"`
	Author              UserIDView `json:"author"`
	Text                string     `json:"text"`
	RawText             string     `json:"rawText"`
	LockedText          bool       `json:"lockedText"`
	IsFavorite          bool       `json:"isFavorite"`
	CanReport           bool       `json:"canReport"`
	CanDelete           bool       `json:"canDelete"`
	CanComment          bool       `json:"canComment"`
	CanEdit             bool       `json:"canEdit"`
	IsPinned            bool       `json:"isPinned"`
	FavoritesCount      int        `json:"favoritesCount"`
	MediaCount          int        `json:"mediaCount"`
	IsMediaReady        bool       `json:"isMediaReady"`
	Voting              any        `json:"voting"`
	IsOpened            bool       `json:"isOpened"`
	CanToggleFavorite   bool       `json:"canToggleFavorite"`
	StreamID            string     `json:"streamId"`
	Price               any        `json:"price"`
	HasVoting           bool       `json:"hasVoting"`
	IsAddedToBookmarks  bool       `json:"isAddedToBookmarks"`
	IsMarkdownDisabled  bool       `json:"isMarkdownDisabled"`
	IsArchived          bool       `json:"isArchived"`
	IsPrivateArchived   bool       `json:"isPrivateArchived"`
	IsDeleted           bool       `json:"isDeleted"`
	HasURL              bool       `json:"hasUrl"`
	IsCouplePeopleMedia bool       `json:"isCouplePeopleMedia"`
	CommentsCount       int        `json:"commentsCount"`
	MentionedUsers      []any      `json:"mentionedUsers"`
	LinkedUsers         []any      `json:"linkedUsers"`
	TipsAmount          string     `json:"tipsAmount"`
	TipsAmountRaw       string     `json:"tipsAmountRaw"`
	Media               []Media    `json:"media"`
	CanViewMedia        bool       `json:"canViewMedia"`
	Preview             []any      `json:"preview"`
}

const (
	MediaTypePhoto = "photo"
	MediaTypeVideo = "video"
	MediaTypeAudio = "audio"
	MediaTypeGif   = "gif"
)

type Media struct {
	ID               int64     `json:"id"`
	Type             string    `json:"type"`
	ConvertedToVideo bool      `json:"convertedToVideo"`
	CanView          bool      `json:"canView"`
	HasError         bool      `json:"hasError"`
	CreatedAt        time.Time `json:"createdAt"`
	Info             struct {
		Source  Source   `json:"source"`
		Preview FileInfo `json:"preview"`
	} `json:"info"`
	Source           Source `json:"source"`
	SquarePreview    string `json:"squarePreview"`
	Full             string `json:"full"`
	Preview          string `json:"preview"`
	Thumb            string `json:"thumb"`
	HasCustomPreview bool   `json:"hasCustomPreview"`
	Duration         int    `json:"duration"`
	IsReady          bool   `json:"isReady"`
	Files            *Files `json:"files"`
	VideoSources     struct {
		Quality720 any `json:"720"`
		Quality240 any `json:"240"`
	} `json:"videoSources"`
}

type Files struct {
	Full          FileInfo `json:"full"`
	Thumb         FileInfo `json:"thumb"`
	Preview       FileInfo `json:"preview"`
	SquarePreview FileInfo `json:"squarePreview"`
	Drm           Drm      `json:"drm"`
}

type FileInfo struct {
	URL     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Size    int    `json:"size"`
	Sources []any  `json:"sources"`
}

type Source struct {
	Source   string `json:"source"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Duration int    `json:"duration"`
}

type Drm struct {
	Manifest  Manifest  `json:"manifest"`
	Signature Signature `json:"signature"`
}

type Manifest struct {
	Hls  string `json:"hls"`
	Dash string `json:"dash"`
}

type Signature struct {
	Hls  CloudFront `json:"hls"`
	Dash CloudFront `json:"dash"`
}

type CloudFront struct {
	CloudFrontPolicy    string `json:"CloudFront-Policy"`
	CloudFrontSignature string `json:"CloudFront-Signature"`
	CloudFrontKeyPairID string `json:"CloudFront-Key-Pair-Id"`
}

type UserIDView struct {
	ID   int64  `json:"id"`
	View string `json:"_view"`
}
