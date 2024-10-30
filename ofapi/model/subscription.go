package model

import "time"

type Subscription struct {
	Avatar       string `json:"avatar"`
	AvatarThumbs struct {
		C144 string `json:"c144"`
		C50  string `json:"c50"`
	} `json:"avatarThumbs"`
	CanAddSubscriber      bool   `json:"canAddSubscriber"`
	CanCommentStory       bool   `json:"canCommentStory"`
	CanEarn               bool   `json:"canEarn"`
	CanLookStory          bool   `json:"canLookStory"`
	CanPayInternal        bool   `json:"canPayInternal"`
	CanReceiveChatMessage bool   `json:"canReceiveChatMessage"`
	CanReport             bool   `json:"canReport"`
	CanRestrict           bool   `json:"canRestrict"`
	CanTrialSend          bool   `json:"canTrialSend"`
	CanUnsubscribe        bool   `json:"canUnsubscribe"`
	CurrentSubscribePrice int    `json:"currentSubscribePrice"`
	DisplayName           string `json:"displayName"`
	HasNotViewedStory     bool   `json:"hasNotViewedStory"`
	HasScheduledStream    bool   `json:"hasScheduledStream"`
	HasStories            bool   `json:"hasStories"`
	HasStream             bool   `json:"hasStream"`
	Header                string `json:"header"`
	HeaderSize            struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"headerSize"`
	HeaderThumbs struct {
		W480 string `json:"w480"`
		W760 string `json:"w760"`
	} `json:"headerThumbs"`
	HideChat          bool      `json:"hideChat"`
	ID                int       `json:"id"`
	IsBlocked         bool      `json:"isBlocked"`
	IsPaywallRequired bool      `json:"isPaywallRequired"`
	IsPerformer       bool      `json:"isPerformer"`
	IsRealPerformer   bool      `json:"isRealPerformer"`
	IsRestricted      bool      `json:"isRestricted"`
	IsVerified        bool      `json:"isVerified"`
	LastSeen          time.Time `json:"lastSeen"`
	ListsStates       []struct {
		CanAddUser          bool   `json:"canAddUser"`
		CannotAddUserReason string `json:"cannotAddUserReason"`
		HasUser             bool   `json:"hasUser"`
		ID                  string `json:"id"`
		Name                string `json:"name"`
		Type                string `json:"type"`
	} `json:"listsStates"`
	Name                    string `json:"name"`
	Notice                  string `json:"notice"`
	SubscribePrice          int    `json:"subscribePrice"`
	SubscribedBy            bool   `json:"subscribedBy"`
	SubscribedByAutoprolong bool   `json:"subscribedByAutoprolong"`
	SubscribedByData        struct {
		DiscountFinishedAt         interface{} `json:"discountFinishedAt"`
		DiscountPercent            int         `json:"discountPercent"`
		DiscountPeriod             int         `json:"discountPeriod"`
		DiscountStartedAt          interface{} `json:"discountStartedAt"`
		Duration                   string      `json:"duration"`
		ExpiredAt                  time.Time   `json:"expiredAt"`
		HasActivePaidSubscriptions bool        `json:"hasActivePaidSubscriptions"`
		IsMuted                    bool        `json:"isMuted"`
		NewPrice                   int         `json:"newPrice"`
		Price                      int         `json:"price"`
		RegularPrice               int         `json:"regularPrice"`
		RenewedAt                  time.Time   `json:"renewedAt"`
		ShowPostsInFeed            bool        `json:"showPostsInFeed"`
		Status                     interface{} `json:"status"`
		SubscribeAt                time.Time   `json:"subscribeAt"`
		SubscribePrice             int         `json:"subscribePrice"`
		Subscribes                 []Subscribe `json:"subscribes"`
		UnsubscribeReason          string      `json:"unsubscribeReason"`
	} `json:"subscribedByData"`
	SubscribedByExpire     bool      `json:"subscribedByExpire"`
	SubscribedByExpireDate time.Time `json:"subscribedByExpireDate"`
	SubscribedIsExpiredNow bool      `json:"subscribedIsExpiredNow"`
	SubscribedOn           bool      `json:"subscribedOn"`
	SubscribedOnData       struct {
		DiscountFinishedAt         interface{} `json:"discountFinishedAt"`
		DiscountPercent            int         `json:"discountPercent"`
		DiscountPeriod             int         `json:"discountPeriod"`
		DiscountStartedAt          interface{} `json:"discountStartedAt"`
		Duration                   string      `json:"duration"`
		ExpiredAt                  time.Time   `json:"expiredAt"`
		HasActivePaidSubscriptions bool        `json:"hasActivePaidSubscriptions"`
		IsMuted                    bool        `json:"isMuted"`
		MessagesSumm               int         `json:"messagesSumm"`
		NewPrice                   int         `json:"newPrice"`
		PostsSumm                  int         `json:"postsSumm"`
		Price                      int         `json:"price"`
		RegularPrice               int         `json:"regularPrice"`
		RenewedAt                  time.Time   `json:"renewedAt"`
		Status                     interface{} `json:"status"`
		StreamsSumm                int         `json:"streamsSumm"`
		SubscribeAt                time.Time   `json:"subscribeAt"`
		SubscribePrice             int         `json:"subscribePrice"`
		Subscribes                 []Subscribe `json:"subscribes"`
		SubscribesSumm             int         `json:"subscribesSumm"`
		TipsSumm                   int         `json:"tipsSumm"`
		TotalSumm                  int         `json:"totalSumm"`
		UnsubscribeReason          string      `json:"unsubscribeReason"`
	} `json:"subscribedOnData"`
	SubscribedOnDuration   string `json:"subscribedOnDuration"`
	SubscribedOnExpiredNow bool   `json:"subscribedOnExpiredNow"`
	TipsEnabled            bool   `json:"tipsEnabled"`
	TipsMax                int    `json:"tipsMax"`
	TipsMin                int    `json:"tipsMin"`
	TipsMinInternal        int    `json:"tipsMinInternal"`
	TipsTextEnabled        bool   `json:"tipsTextEnabled"`
	Username               string `json:"username"`
	View                   string `json:"view"`
}

type Subscribe struct {
	Action       string      `json:"action"`
	CancelDate   interface{} `json:"cancelDate"`
	Date         time.Time   `json:"date"`
	Discount     int         `json:"discount"`
	Duration     int         `json:"duration"`
	EarningID    int         `json:"earningId"`
	ExpireDate   time.Time   `json:"expireDate"`
	ID           int         `json:"id"`
	IsCurrent    bool        `json:"isCurrent"`
	OfferEnd     interface{} `json:"offerEnd"`
	OfferStart   interface{} `json:"offerStart"`
	Price        int         `json:"price"`
	RegularPrice int         `json:"regularPrice"`
	StartDate    time.Time   `json:"startDate"`
	SubscriberID int         `json:"subscriberId"`
	Type         string      `json:"type"`
	UserID       int         `json:"userId"`
}
