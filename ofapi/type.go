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
