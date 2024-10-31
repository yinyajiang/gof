package model

import "time"

type Subscription struct {
	View                    string         `json:"view"`
	Avatar                  string         `json:"avatar"`
	AvatarThumbs            AvatarThumbs   `json:"avatarThumbs"`
	Header                  string         `json:"header"`
	HeaderSize              HeaderSize     `json:"headerSize"`
	HeaderThumbs            HeaderThumbs   `json:"headerThumbs"`
	ID                      int64          `json:"id"`
	Name                    string         `json:"name"`
	Username                string         `json:"username"`
	CanLookStory            bool           `json:"canLookStory"`
	CanCommentStory         bool           `json:"canCommentStory"`
	HasNotViewedStory       bool           `json:"hasNotViewedStory"`
	IsVerified              bool           `json:"isVerified"`
	CanPayInternal          bool           `json:"canPayInternal"`
	HasScheduledStream      bool           `json:"hasScheduledStream"`
	HasStream               bool           `json:"hasStream"`
	HasStories              bool           `json:"hasStories"`
	TipsEnabled             bool           `json:"tipsEnabled"`
	TipsTextEnabled         bool           `json:"tipsTextEnabled"`
	TipsMin                 int            `json:"tipsMin"`
	TipsMinInternal         int            `json:"tipsMinInternal"`
	TipsMax                 int            `json:"tipsMax"`
	CanEarn                 bool           `json:"canEarn"`
	CanAddSubscriber        bool           `json:"canAddSubscriber"`
	SubscribePrice          any            `json:"subscribePrice"`
	IsPaywallRequired       bool           `json:"isPaywallRequired"`
	Unprofitable            bool           `json:"unprofitable"`
	ListsStates             []ListsState   `json:"listsStates"`
	IsMuted                 bool           `json:"isMuted"`
	IsRestricted            bool           `json:"isRestricted"`
	CanRestrict             bool           `json:"canRestrict"`
	SubscribedBy            bool           `json:"subscribedBy"`
	SubscribedByExpire      bool           `json:"subscribedByExpire"`
	SubscribedByExpireDate  time.Time      `json:"subscribedByExpireDate"`
	SubscribedByAutoprolong bool           `json:"subscribedByAutoprolong"`
	SubscribedIsExpiredNow  bool           `json:"subscribedIsExpiredNow"`
	CurrentSubscribePrice   any            `json:"currentSubscribePrice"`
	SubscribedOn            bool           `json:"subscribedOn"`
	SubscribedOnExpiredNow  bool           `json:"subscribedOnExpiredNow"`
	SubscribedOnDuration    string         `json:"subscribedOnDuration"`
	CanReport               bool           `json:"canReport"`
	CanReceiveChatMessage   bool           `json:"canReceiveChatMessage"`
	HideChat                bool           `json:"hideChat"`
	LastSeen                time.Time      `json:"lastSeen"`
	IsPerformer             bool           `json:"isPerformer"`
	IsRealPerformer         bool           `json:"isRealPerformer"`
	SubscribedByData        SubscribedData `json:"subscribedByData"`
	SubscribedOnData        SubscribedData `json:"subscribedOnData"`
	CanTrialSend            bool           `json:"canTrialSend"`
	IsBlocked               bool           `json:"isBlocked"`
	DisplayName             string         `json:"displayName"`
	Notice                  string         `json:"notice"`
}

type ListsState struct {
	ID         any    `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	HasUser    bool   `json:"hasUser"`
	CanAddUser bool   `json:"canAddUser"`
}

type Subscribe struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"userId"`
	SubscriberID int64     `json:"subscriberId"`
	Date         time.Time `json:"date"`
	Duration     int       `json:"duration"`
	StartDate    time.Time `json:"startDate"`
	ExpireDate   time.Time `json:"expireDate"`
	CancelDate   any       `json:"cancelDate"`
	Price        any       `json:"price"`
	RegularPrice any       `json:"regularPrice"`
	Discount     int       `json:"discount"`
	EarningID    int       `json:"earningId"`
	Action       string    `json:"action"`
	Type         string    `json:"type"`
	OfferStart   any       `json:"offerStart"`
	OfferEnd     any       `json:"offerEnd"`
	IsCurrent    bool      `json:"isCurrent"`
}

type SubscribedData struct {
	Price                      any         `json:"price"`
	NewPrice                   any         `json:"newPrice"`
	RegularPrice               any         `json:"regularPrice"`
	SubscribePrice             any         `json:"subscribePrice"`
	DiscountPercent            int         `json:"discountPercent"`
	DiscountPeriod             int         `json:"discountPeriod"`
	SubscribeAt                time.Time   `json:"subscribeAt"`
	ExpiredAt                  time.Time   `json:"expiredAt"`
	RenewedAt                  time.Time   `json:"renewedAt"`
	DiscountFinishedAt         any         `json:"discountFinishedAt"`
	DiscountStartedAt          any         `json:"discountStartedAt"`
	Status                     any         `json:"status"`
	IsMuted                    bool        `json:"isMuted"`
	UnsubscribeReason          any         `json:"unsubscribeReason"`
	Duration                   any         `json:"duration"`
	TipsSumm                   any         `json:"tipsSumm"`
	SubscribesSumm             any         `json:"subscribesSumm"`
	MessagesSumm               any         `json:"messagesSumm"`
	PostsSumm                  any         `json:"postsSumm"`
	StreamsSumm                any         `json:"streamsSumm"`
	TotalSumm                  any         `json:"totalSumm"`
	Subscribes                 []Subscribe `json:"subscribes"`
	ShowPostsInFeed            bool        `json:"showPostsInFeed"`
	HasActivePaidSubscriptions bool        `json:"hasActivePaidSubscriptions"`
	LastActivity               time.Time   `json:"lastActivity"`
}

type SubscriptionBundle struct {
	ID       int64  `json:"id"`
	Discount any    `json:"discount"`
	Duration string `json:"duration"`
	Price    any    `json:"price"`
	CanBuy   bool   `json:"canBuy"`
}
