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
	AppToken         string `json:"app_token"`
	ChecksumConstant int    `json:"checksum_constant"`
	ChecksumIndexes  []int  `json:"checksum_indexes"`
	Prefix           string `json:"prefix"`
	StaticParam      string `json:"static_param"`
	Suffix           string `json:"suffix"`
	Revision         string `json:"revision"`
}
