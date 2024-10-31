package model

import "time"

type User struct {
	AgeVerificationRequired   any            `json:"ageVerificationRequired"`
	View                      string         `json:"view"`
	Avatar                    string         `json:"avatar"`
	AvatarThumbs              AvatarThumbs   `json:"avatarThumbs"`
	Header                    any            `json:"header"`
	HeaderSize                HeaderSize     `json:"headerSize"`
	HeaderThumbs              HeaderThumbs   `json:"headerThumbs"`
	ID                        int64          `json:"id"`
	Name                      string         `json:"name"`
	Username                  string         `json:"username"`
	CanLookStory              any            `json:"canLookStory"`
	CanCommentStory           bool           `json:"canCommentStory"`
	HasNotViewedStory         bool           `json:"hasNotViewedStory"`
	IsVerified                bool           `json:"isVerified"`
	CanPayInternal            bool           `json:"canPayInternal"`
	HasScheduledStream        bool           `json:"hasScheduledStream"`
	HasStream                 bool           `json:"hasStream"`
	HasStories                bool           `json:"hasStories"`
	TipsEnabled               bool           `json:"tipsEnabled"`
	TipsTextEnabled           bool           `json:"tipsTextEnabled"`
	TipsMin                   int            `json:"tipsMin"`
	TipsMinInternal           int            `json:"tipsMinInternal"`
	TipsMax                   int            `json:"tipsMax"`
	CanEarn                   bool           `json:"canEarn"`
	CanAddSubscriber          bool           `json:"canAddSubscriber"`
	SubscribePrice            any            `json:"subscribePrice"`
	DisplayName               string         `json:"displayName"`
	Notice                    any            `json:"notice"`
	IsPaywallRequired         bool           `json:"isPaywallRequired"`
	Unprofitable              bool           `json:"unprofitable"`
	ListsStates               []ListsState   `json:"listsStates"`
	IsMuted                   bool           `json:"isMuted"`
	IsRestricted              bool           `json:"isRestricted"`
	CanRestrict               bool           `json:"canRestrict"`
	SubscribedBy              any            `json:"subscribedBy"`
	SubscribedByExpire        any            `json:"subscribedByExpire"`
	SubscribedByExpireDate    time.Time      `json:"subscribedByExpireDate"`
	SubscribedByAutoprolong   any            `json:"subscribedByAutoprolong"`
	SubscribedIsExpiredNow    bool           `json:"subscribedIsExpiredNow"`
	CurrentSubscribePrice     any            `json:"currentSubscribePrice"`
	SubscribedOn              any            `json:"subscribedOn"`
	SubscribedOnExpiredNow    any            `json:"subscribedOnExpiredNow"`
	SubscribedOnDuration      any            `json:"subscribedOnDuration"`
	JoinDate                  time.Time      `json:"joinDate"`
	IsReferrerAllowed         bool           `json:"isReferrerAllowed"`
	About                     any            `json:"about"`
	RawAbout                  any            `json:"rawAbout"`
	WsAuthToken               string         `json:"wsAuthToken"`
	WsURL                     string         `json:"wsUrl"`
	Website                   any            `json:"website"`
	Wishlist                  any            `json:"wishlist"`
	Location                  any            `json:"location"`
	PostsCount                int            `json:"postsCount"`
	ArchivedPostsCount        int            `json:"archivedPostsCount"`
	PrivateArchivedPostsCount int            `json:"privateArchivedPostsCount"`
	PhotosCount               int            `json:"photosCount"`
	VideosCount               int            `json:"videosCount"`
	AudiosCount               int            `json:"audiosCount"`
	MediasCount               int            `json:"mediasCount"`
	LastSeen                  any            `json:"lastSeen"`
	FavoritesCount            int            `json:"favoritesCount"`
	FavoritedCount            int            `json:"favoritedCount"`
	ShowPostsInFeed           bool           `json:"showPostsInFeed"`
	CanReceiveChatMessage     bool           `json:"canReceiveChatMessage"`
	IsPerformer               bool           `json:"isPerformer"`
	IsRealPerformer           bool           `json:"isRealPerformer"`
	IsSpotifyConnected        bool           `json:"isSpotifyConnected"`
	SubscribersCount          int            `json:"subscribersCount"`
	HasPinnedPosts            bool           `json:"hasPinnedPosts"`
	HasLabels                 bool           `json:"hasLabels"`
	CanChat                   bool           `json:"canChat"`
	CallPrice                 any            `json:"callPrice"`
	IsPrivateRestriction      bool           `json:"isPrivateRestriction"`
	ShowSubscribersCount      any            `json:"showSubscribersCount"`
	ShowMediaCount            any            `json:"showMediaCount"`
	SubscribedByData          SubscribedData `json:"subscribedByData"`
	SubscribedOnData          SubscribedData `json:"subscribedOnData"`
	CanPromotion              bool           `json:"canPromotion"`
	CanCreatePromotion        bool           `json:"canCreatePromotion"`
	CanCreateTrial            bool           `json:"canCreateTrial"`
	IsAdultContent            bool           `json:"isAdultContent"`
	CanTrialSend              bool           `json:"canTrialSend"`
	HadEnoughLastPhotos       bool           `json:"hadEnoughLastPhotos"`
	HasLinks                  bool           `json:"hasLinks"`
	FirstPublishedPostDate    time.Time      `json:"firstPublishedPostDate"`
	IsSpringConnected         bool           `json:"isSpringConnected"`
	IsFriend                  bool           `json:"isFriend"`
	IsBlocked                 bool           `json:"isBlocked"`
	CanReport                 bool           `json:"canReport"`
	CanAddCard                bool           `json:"canAddCard"`
	Email                     string         `json:"email"`
	IP                        any            `json:"ip"`
	IsAuth                    bool           `json:"isAuth"`
}

type AvatarThumbs struct {
	C50  string `json:"c50"`
	C144 string `json:"c144"`
}

type HeaderSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type HeaderThumbs struct {
	W480 string `json:"w480"`
	W760 string `json:"w760"`
}
