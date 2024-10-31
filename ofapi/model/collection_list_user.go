package model

import "time"

type CollectionListUser struct {
	View                    string               `json:"view"`
	Avatar                  string               `json:"avatar"`
	AvatarThumbs            AvatarThumbs         `json:"avatarThumbs"`
	Header                  string               `json:"header"`
	HeaderSize              HeaderSize           `json:"headerSize"`
	HeaderThumbs            HeaderThumbs         `json:"headerThumbs"`
	ID                      int64                `json:"id"`
	Name                    string               `json:"name"`
	Username                string               `json:"username"`
	CanLookStory            bool                 `json:"canLookStory"`
	CanCommentStory         bool                 `json:"canCommentStory"`
	HasNotViewedStory       bool                 `json:"hasNotViewedStory"`
	IsVerified              bool                 `json:"isVerified"`
	CanPayInternal          bool                 `json:"canPayInternal"`
	HasScheduledStream      bool                 `json:"hasScheduledStream"`
	HasStream               bool                 `json:"hasStream"`
	HasStories              bool                 `json:"hasStories"`
	TipsEnabled             bool                 `json:"tipsEnabled"`
	TipsTextEnabled         bool                 `json:"tipsTextEnabled"`
	TipsMin                 int                  `json:"tipsMin"`
	TipsMinInternal         int                  `json:"tipsMinInternal"`
	TipsMax                 int                  `json:"tipsMax"`
	CanEarn                 bool                 `json:"canEarn"`
	CanAddSubscriber        bool                 `json:"canAddSubscriber"`
	SubscribePrice          any                  `json:"subscribePrice"`
	SubscriptionBundles     []SubscriptionBundle `json:"subscriptionBundles"`
	DisplayName             string               `json:"displayName"`
	Notice                  string               `json:"notice"`
	IsPaywallRequired       bool                 `json:"isPaywallRequired"`
	Unprofitable            bool                 `json:"unprofitable"`
	ListsStates             []ListsState         `json:"listsStates"`
	IsMuted                 bool                 `json:"isMuted"`
	IsRestricted            bool                 `json:"isRestricted"`
	CanRestrict             bool                 `json:"canRestrict"`
	SubscribedBy            bool                 `json:"subscribedBy"`
	SubscribedByExpire      bool                 `json:"subscribedByExpire"`
	SubscribedByExpireDate  time.Time            `json:"subscribedByExpireDate"`
	SubscribedByAutoprolong bool                 `json:"subscribedByAutoprolong"`
	SubscribedIsExpiredNow  bool                 `json:"subscribedIsExpiredNow"`
	CurrentSubscribePrice   any                  `json:"currentSubscribePrice"`
	SubscribedOn            bool                 `json:"subscribedOn"`
	SubscribedOnExpiredNow  bool                 `json:"subscribedOnExpiredNow"`
	SubscribedOnDuration    string               `json:"subscribedOnDuration"`
	CanReport               bool                 `json:"canReport"`
	CanReceiveChatMessage   bool                 `json:"canReceiveChatMessage"`
	HideChat                bool                 `json:"hideChat"`
	LastSeen                time.Time            `json:"lastSeen"`
	IsPerformer             bool                 `json:"isPerformer"`
	IsRealPerformer         bool                 `json:"isRealPerformer"`
	SubscribedByData        SubscribedData       `json:"subscribedByData"`
	SubscribedOnData        SubscribedData       `json:"subscribedOnData"`
	CanTrialSend            bool                 `json:"canTrialSend"`
	IsBlocked               bool                 `json:"isBlocked"`
	PromoOffers             []any                `json:"promoOffers"`
}
