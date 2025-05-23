package ofapi

import (
	"reflect"

	"github.com/yinyajiang/gof/ofapi/model"
)

type SubscritionType string

const (
	SubscritionTypeActive  SubscritionType = "active"
	SubscritionTypeExpired SubscritionType = "expired"
	SubscritionTypeAll     SubscritionType = "all"
)

type UserMedias string

const (
	UserVideos UserMedias = "videos"
	UserPhotos UserMedias = "photos"
	UserAll    UserMedias = "all"
)

type BookmarkMedia string

const (
	BookmarkPhotos BookmarkMedia = "photos"
	BookmarkVideos BookmarkMedia = "videos"
	BookmarkAudios BookmarkMedia = "audios"
	BookmarkOther  BookmarkMedia = "other"
	BookmarkLocked BookmarkMedia = "locked"
	BookmarkAll    BookmarkMedia = "all"
)

type SubscribeFilter func(sub model.Subscription) bool

func SubscribeRestrictedFilter(includeRestricted bool) SubscribeFilter {
	return func(sub model.Subscription) bool {
		return includeRestricted || !sub.IsRestricted || (sub.IsRestricted && includeRestricted)
	}
}

type CollectionFilter func(collection model.Collection) bool

func CollectionFilterByID(id any) CollectionFilter {
	return func(collection model.Collection) bool {
		return reflect.DeepEqual(collection.ID, id)
	}
}

const (
	CollectionTypeFans      = "fans"
	CollectionTypeFollowing = "following"
	CollectionTypeCustom    = "custom"
)

func CollectionFilterByType(collectionType string) CollectionFilter {
	return func(collection model.Collection) bool {
		return collection.Type == collectionType
	}
}

type TimeDirection int

const (
	TimeDirectionBefore TimeDirection = 0
	TimeDirectionAfter  TimeDirection = 1
)

type UserIdentifier struct {
	ID       int64
	Username string
}

type rules struct {
	AppToken         string `json:"app-token"`
	AppToken_Old     string `json:"app_token"` //old config
	ChecksumConstant int    `json:"checksum_constant"`
	ChecksumIndexes  []int  `json:"checksum_indexes"`
	Prefix           string `json:"prefix"`
	StaticParam      string `json:"static_param"`
	Suffix           string `json:"suffix"`
	Revision         string `json:"revision"`
}

type OFAuthInfo struct {
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	X_BC      string `json:"x_bc"`
	Cookie    string `json:"cookie"`
}

func (authInfo OFAuthInfo) String() string {
	return authInfo2String(authInfo)
}

func (authInfo OFAuthInfo) IsEmpty() bool {
	return authInfo.UserID == "" ||
		authInfo.UserAgent == "" ||
		authInfo.X_BC == "" ||
		authInfo.Cookie == ""
}
