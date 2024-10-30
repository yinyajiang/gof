package ofapi

import "github.com/yinyajiang/gof/ofapi/model"

type SubscritionType string

const (
	SubscritionTypeActive  SubscritionType = "active"
	SubscritionTypeExpired SubscritionType = "expired"
)

type SubscribeFilter func(sub model.Subscription) bool

func SubscribeRestrictedFilter(includeRestricted bool) SubscribeFilter {
	return func(sub model.Subscription) bool {
		return includeRestricted || !sub.IsRestricted || (sub.IsRestricted && includeRestricted)
	}
}
