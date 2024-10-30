package model

import "time"

type UserInfo struct {
	About                   string   `json:"about"`
	AdvBlock                []string `json:"advBlock"`
	AgeVerificationRequired bool     `json:"ageVerificationRequired"`
	AgeVerificationSession  struct {
		APIFlow   string    `json:"apiFlow"`
		ExpiredAt time.Time `json:"expiredAt"`
		Status    string    `json:"status"`
		URL       string    `json:"url"`
	} `json:"ageVerificationSession"`
	ArchivedPostsCount int    `json:"archivedPostsCount"`
	AudiosCount        int    `json:"audiosCount"`
	Avatar             string `json:"avatar"`
	AvatarThumbs       struct {
		C50  string `json:"c50"`
		C144 string `json:"c144"`
	} `json:"avatarThumbs"`
	CanAddCard                bool          `json:"canAddCard"`
	CanAlternativeWalletTopUp bool          `json:"canAlternativeWalletTopUp"`
	CanChat                   bool          `json:"canChat"`
	CanCommentStory           bool          `json:"canCommentStory"`
	CanConnectOfAccount       bool          `json:"canConnectOfAccount"`
	CanCreateLists            bool          `json:"canCreateLists"`
	CanLookStory              bool          `json:"canLookStory"`
	CanPayInternal            bool          `json:"canPayInternal"`
	CanPinPost                bool          `json:"canPinPost"`
	CanReceiveChatMessage     bool          `json:"canReceiveChatMessage"`
	CanSendChatToAll          bool          `json:"canSendChatToAll"`
	ChatMessagesCount         int           `json:"chatMessagesCount"`
	ConnectedOfAccounts       []interface{} `json:"connectedOfAccounts"`
	CountPinnedChat           int           `json:"countPinnedChat"`
	CountPriorityChat         int           `json:"countPriorityChat"`
	CreditBalance             int           `json:"creditBalance"`
	CreditsMax                int           `json:"creditsMax"`
	CreditsMin                int           `json:"creditsMin"`
	Csrf                      string        `json:"csrf"`
	Email                     string        `json:"email"`
	EnabledImageEditorForChat bool          `json:"enabledImageEditorForChat"`
	FaceIDRegular             struct {
		CanPostpone bool `json:"canPostpone"`
		NeedToShow  bool `json:"needToShow"`
	} `json:"faceIdRegular"`
	FavoritedCount                  int  `json:"favoritedCount"`
	FavoritesCount                  int  `json:"favoritesCount"`
	HasInternalPayments             bool `json:"hasInternalPayments"`
	HasLabels                       bool `json:"hasLabels"`
	HasNewAlerts                    bool `json:"hasNewAlerts"`
	HasNewChangedPriceSubscriptions bool `json:"hasNewChangedPriceSubscriptions"`
	HasNewHints                     bool `json:"hasNewHints"`
	HasNewTicketReplies             struct {
		Closed bool `json:"closed"`
		Open   bool `json:"open"`
		Solved bool `json:"solved"`
	} `json:"hasNewTicketReplies"`
	HasNotViewedStory      bool   `json:"hasNotViewedStory"`
	HasPinnedPosts         bool   `json:"hasPinnedPosts"`
	HasPurchasedPosts      bool   `json:"hasPurchasedPosts"`
	HasScenario            bool   `json:"hasScenario"`
	HasSystemNotifications bool   `json:"hasSystemNotifications"`
	HasTags                bool   `json:"hasTags"`
	HasWatermarkPhoto      bool   `json:"hasWatermarkPhoto"`
	HasWatermarkVideo      bool   `json:"hasWatermarkVideo"`
	Header                 string `json:"header"`
	HeaderSize             struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"headerSize"`
	HeaderThumbs struct {
		W480 string `json:"w480"`
		W760 string `json:"w760"`
	} `json:"headerThumbs"`
	ID                        int         `json:"id"`
	IP                        string      `json:"ip"`
	IsAgeVerified             bool        `json:"isAgeVerified"`
	IsAllowTweets             bool        `json:"isAllowTweets"`
	IsAuth                    bool        `json:"isAuth"`
	IsCreditsEnabled          bool        `json:"isCreditsEnabled"`
	IsDeleteInitiated         bool        `json:"isDeleteInitiated"`
	IsEmailChecked            bool        `json:"isEmailChecked"`
	IsEmailRequired           bool        `json:"isEmailRequired"`
	IsLegalApprovedAllowed    bool        `json:"isLegalApprovedAllowed"`
	IsMakePayment             bool        `json:"isMakePayment"`
	IsOtpEnabled              bool        `json:"isOtpEnabled"`
	IsPaymentCardConnected    bool        `json:"isPaymentCardConnected"`
	IsPaywallPassed           bool        `json:"isPaywallPassed"`
	IsPerformer               bool        `json:"isPerformer"`
	IsRealCardConnected       bool        `json:"isRealCardConnected"`
	IsRealPerformer           bool        `json:"isRealPerformer"`
	IsReferrerAllowed         bool        `json:"isReferrerAllowed"`
	IsSpotifyConnected        bool        `json:"isSpotifyConnected"`
	IsTwitterConnected        bool        `json:"isTwitterConnected"`
	IsVerified                bool        `json:"isVerified"`
	IsVisibleOnline           bool        `json:"isVisibleOnline"`
	IsWalletAutorecharge      bool        `json:"isWalletAutorecharge"`
	IsWantComments            bool        `json:"isWantComments"`
	IvFlow                    string      `json:"ivFlow"`
	JoinDate                  time.Time   `json:"joinDate"`
	LastSeen                  time.Time   `json:"lastSeen"`
	Location                  interface{} `json:"location"`
	MaxPinnedPostsCount       int         `json:"maxPinnedPostsCount"`
	MediasCount               int         `json:"mediasCount"`
	Name                      string      `json:"name"`
	NeedIVApprove             bool        `json:"needIVApprove"`
	NewTagsCount              int         `json:"newTagsCount"`
	NotificationsCount        int         `json:"notificationsCount"`
	PaidFeed                  bool        `json:"paidFeed"`
	PayoutLegalApproveState   string      `json:"payoutLegalApproveState"`
	PhotosCount               int         `json:"photosCount"`
	PinnedPostsCount          int         `json:"pinnedPostsCount"`
	PostsCount                int         `json:"postsCount"`
	PrivateArchivedPostsCount int         `json:"privateArchivedPostsCount"`
	ShowPostsInFeed           bool        `json:"showPostsInFeed"`
	SubscribersCount          int         `json:"subscribersCount"`
	SubscribesCount           int         `json:"subscribesCount"`
	TwitterUsername           interface{} `json:"twitterUsername"`
	Upload                    struct {
		GeoUploadArgs struct {
			Additional struct {
				User string `json:"user"`
			} `json:"additional"`
			IsDelay         bool   `json:"isDelay"`
			NeedThumbs      bool   `json:"needThumbs"`
			Preset          string `json:"preset"`
			PresetPng       string `json:"preset_png"`
			ProtectedPreset string `json:"protected_preset"`
		} `json:"geoUploadArgs"`
	} `json:"upload"`
	Username                 string      `json:"username"`
	VideosCount              int         `json:"videosCount"`
	View                     string      `json:"view"`
	WalletAutorechargeAmount int         `json:"walletAutorechargeAmount"`
	WalletAutorechargeMin    int         `json:"walletAutorechargeMin"`
	WalletFirstRebills       bool        `json:"walletFirstRebills"`
	WatermarkPosition        string      `json:"watermarkPosition"`
	WatermarkText            string      `json:"watermarkText"`
	Website                  interface{} `json:"website"`
	Wishlist                 interface{} `json:"wishlist"`
	WsAuthToken              string      `json:"wsAuthToken"`
	WsURL                    string      `json:"wsUrl"`
}
