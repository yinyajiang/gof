package ofapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi/model"
	"github.com/yinyajiang/gof/ofwebview"
)

var ErrorAuth = errors.New("need auth")
var ErrorChangeAccount = fmt.Errorf("%w: %s", ErrorAuth, "try change login account")

type OFAPI struct {
	req *Req

	cacheDir string
}

type OFApiConfig struct {
	OptionalRulesURI []string
	ApiCacheDir      string
}

func NewOFAPI(config OFApiConfig) (*OFAPI, error) {
	api := &OFAPI{
		req:      &Req{},
		cacheDir: config.ApiCacheDir,
	}

	rules, err := loadRules(api.cacheDir, config.OptionalRulesURI)
	if err != nil {
		return nil, err
	}
	api.req.SetRules(rules)

	api.AuthByCache()
	return api, nil
}

func (c *OFAPI) Req() *Req {
	return c.req
}

func (c *OFAPI) HasAuthInfo() bool {
	authInfo := c.req.AuthInfo()
	rules := c.req.Rules()
	return authInfo.Cookie != "" &&
		authInfo.X_BC != "" &&
		authInfo.UserAgent != "" &&
		rules.AppToken != ""
}

/*
user_id:={} || user_agent:={} || x_bc:={} || cookie:={ sess={};auth_id={} }
*/
func (c *OFAPI) AuthByString(authInfo string, check ...bool) error {
	if authInfo == "" {
		return errors.New("authInfo is empty")
	}
	return c.Auth(String2AuthInfo(authInfo), check...)
}

func (c *OFAPI) AuthByWebviewLoginResult(info ofwebview.LoginResult, check ...bool) error {
	authInfo, err := cookies2AuthInfo(info.UA, info.Cookies)
	if err != nil {
		return err
	}
	return c.Auth(authInfo, check...)
}

func (c *OFAPI) AuthByCookieFile(ua, cookiefile string, check ...bool) error {
	info, err := CookieFile2AuthInfo(ua, cookiefile)
	if err != nil {
		return err
	}
	return c.Auth(info, check...)
}

func (c *OFAPI) AuthByCache(check ...bool) error {
	return c.Auth(OFAuthInfo{}, check...)
}

func (c *OFAPI) Auth(authInfo OFAuthInfo, check ...bool) error {
	authInfo = correctAuthInfo(authInfo)

	if c.req.AuthInfo().String() == authInfo.String() && !authInfo.IsEmpty() {
		fmt.Println("authInfo is same, skip")
		return nil
	}

	if gof.IsDebug() {
		fmt.Printf("new authInfo: %s\n", authInfo.String())
		fmt.Printf("old authInfo: %s\n", c.req.AuthInfo().String())
	}

	//from cache
	if authInfo.IsEmpty() {
		if c.HasAuthInfo() {
			return errors.New("AuthInfo is invalid")
		}
		auth, err := loadCacheAuthInfo(c.cacheDir)
		if err != nil {
			return errors.New("AuthInfo is invalid")
		}
		c.req.SetAuthInfo(auth)
		return nil
	}

	if authInfo.Cookie == "" || authInfo.X_BC == "" || authInfo.UserAgent == "" {
		return errors.New("AuthInfo is invalid")
	}

	old := c.req.AuthInfo()
	c.req.SetAuthInfo(authInfo)
	if len(check) != 0 && check[0] {
		err := c.CheckAuth()
		if err != nil {
			//restore
			c.req.SetAuthInfo(old)
			return err
		}
	}
	cacheAuthInfo(c.cacheDir, authInfo)
	return nil
}

func (c *OFAPI) CheckAuth() error {
	if !c.HasAuthInfo() {
		return ErrorAuth
	}

	me, err := c.GetMe()
	if err != nil {
		return fmt.Errorf("%w, GetMe failed: %v", ErrorAuth, err)
	}
	if me.Username == "" && me.Name == "" {
		return ErrorAuth
	}
	return nil
}

func (c *OFAPI) GetMe() (model.User, error) {
	return c.GetUser("me")
}

func (c *OFAPI) UserAgent() string {
	return c.req.UserAgent()
}

func (c *OFAPI) GetChatMessages(userID int64) ([]model.Post, error) {
	var result []model.Post
	var nextID *int64
	var hasMore = true
	var err error
	for hasMore {
		param := map[string]string{
			"limit": "50",
			"order": "desc",
		}
		if nextID != nil {
			param["id"] = strconv.FormatInt(*nextID, 10)
		}

		var moreList moreList[model.Post]
		err = c.req.GetUnmarshal(ApiURLPath("/chats/%d/messages", userID), param, &moreList)
		if err != nil {
			break
		}
		hasMore = moreList.HasMore
		if len(moreList.List) != 0 {
			nextID = &moreList.List[len(moreList.List)-1].ID
			result = append(result, moreList.List...)
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return nil, err
}

func (c *OFAPI) GetUserHightlights(userID int64, withStories bool) ([]model.Highlight, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.Highlight
	for hasMore {
		var moreList moreList[model.Highlight]
		err = c.req.GetUnmarshal(ApiURLPath("/users/%d/stories/highlights", userID), map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "5",
		}, &moreList)
		if err != nil {
			break
		}

		hasMore = moreList.HasMore
		offset += len(moreList.List)

		if withStories {
			for _, highlight := range moreList.List {
				highlight, e := c.GetHighlight(highlight.ID)
				if e == nil {
					result = append(result, highlight)
				}
			}
		} else {
			result = append(result, moreList.List...)
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return nil, err
}

func (c *OFAPI) GetHighlight(highlightID int64) (model.Highlight, error) {
	var highlight model.Highlight
	err := c.req.GetUnmarshal(ApiURLPath("/stories/highlights/%d", highlightID), nil, &highlight)
	return highlight, err
}

func (c *OFAPI) GetUserStories(userID int64) ([]model.Story, error) {
	var stories []model.Story
	err := c.req.GetUnmarshal(ApiURLPath("/users/%d/stories", userID), map[string]string{
		"limit": "50",
		"order": "publish_date_desc",
	}, &stories)
	return stories, err
}

func (c *OFAPI) GetPaidPosts() ([]model.Post, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.Post

	for hasMore {
		var moreList moreList[model.Post]
		err = c.req.GetUnmarshal("/posts/paid", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
			"order":  "publish_date_desc",
			"format": "infinite",
		}, &moreList)
		if err != nil {
			break
		}
		hasMore = moreList.HasMore
		offset += len(moreList.List)

		result = append(result, moreList.List...)
	}
	return result, err
}

func (c *OFAPI) GetPost(postID any) (model.Post, error) {
	var post model.Post
	err := c.req.GetUnmarshal(ApiURLPath("/posts/%v", postID), map[string]string{
		"skip_users": "all",
	}, &post)
	return post, err
}

func (c *OFAPI) GetAllBookmarkes(bookmarkMedia BookmarkMedia) ([]model.Post, error) {
	var endpoint string
	switch bookmarkMedia {
	case BookmarkPhotos, BookmarkVideos, BookmarkAudios, BookmarkOther, BookmarkLocked:
		endpoint = "/" + string(bookmarkMedia)
	default:
		endpoint = "/all"
	}
	return c.getBookmarkesByEndPoint(endpoint)
}

func (c *OFAPI) GetBookmark(bookmarkID any, bookmarkMedia BookmarkMedia) ([]model.Post, error) {
	var endpoint string
	switch bookmarkMedia {
	case BookmarkPhotos, BookmarkVideos, BookmarkAudios, BookmarkOther, BookmarkLocked:
		endpoint = "/" + string(bookmarkMedia)
	default:
		endpoint = "/all"
	}
	if bookmarkID != nil {
		endpoint += fmt.Sprintf("/%v", bookmarkID)
	}

	return c.getBookmarkesByEndPoint(strings.TrimRight(endpoint, "/"))
}

func (c *OFAPI) getBookmarkesByEndPoint(endpoint string) ([]model.Post, error) {
	if endpoint != "" && !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	var err error
	hasMore := true
	offset := 0
	var result []model.Post
	for hasMore {
		var moreList moreList[model.Post]
		err = c.req.GetUnmarshal("/posts/bookmarks"+endpoint, map[string]string{
			"offset":     strconv.Itoa(offset),
			"limit":      "50",
			"format":     "infinite",
			"skip_users": "all",
		}, &moreList)
		if err != nil {
			break
		}
		hasMore = moreList.HasMore
		offset += len(moreList.List)
		result = append(result, moreList.List...)
	}
	if len(result) != 0 {
		return result, nil
	}
	return nil, err
}

func (c *OFAPI) GetUserPosts(userID int64) ([]model.Post, error) {
	return c.GetUserPostsByTime(userID, time.Now(), TimeDirectionBefore)
}

func (c *OFAPI) GetUserMedias(userID int64, userMedias UserMedias) ([]model.Post, error) {
	return c.GetUserMediasByTime(userID, time.Now(), TimeDirectionBefore, userMedias)
}

func (c *OFAPI) GetUserStreams(userID int64) ([]model.Post, error) {
	return c.GetUserStreamsByTime(userID, time.Now(), TimeDirectionBefore)
}

func (c *OFAPI) GetUserArchived(userID int64) ([]model.Post, error) {
	return c.GetUserArchivedByTime(userID, time.Now(), TimeDirectionBefore)
}

func (c *OFAPI) GetUserArchivedByTime(userID int64, timePoint time.Time, timeDirection TimeDirection) ([]model.Post, error) {
	return c.getUserPostsByEndPointAndTime(userID, "/posts", map[string]string{
		"skip_users": "all",
		"label":      "archived",
		"counters":   "1",
	}, timePoint, timeDirection)
}

func (c *OFAPI) GetUserStreamsByTime(userID int64, timePoint time.Time, timeDirection TimeDirection) ([]model.Post, error) {
	return c.getUserPostsByEndPointAndTime(userID, "/posts/streams", nil, timePoint, timeDirection)
}

func (c *OFAPI) GetUserMediasByTime(userID int64, timePoint time.Time, timeDirection TimeDirection, userMedias UserMedias) ([]model.Post, error) {

	var endpoint string
	switch userMedias {
	case UserVideos, UserPhotos:
		endpoint = "/posts/" + string(userMedias)
	default:
		endpoint = "/posts/medias"
	}

	return c.getUserPostsByEndPointAndTime(userID, endpoint, map[string]string{
		"skip_users": "all",
	}, timePoint, timeDirection)
}

func (c *OFAPI) GetUserPostsByTime(userID int64, timePoint time.Time, timeDirection TimeDirection) ([]model.Post, error) {
	return c.getUserPostsByEndPointAndTime(userID, "/posts", nil, timePoint, timeDirection)
}

func (c *OFAPI) getUserPostsByEndPointAndTime(userID int64, endpoint string, mergeParam map[string]string, timePoint time.Time, timeDirection TimeDirection) ([]model.Post, error) {
	param := initPublishTimeParam(map[string]string{
		"limit":  "50",
		"order":  "publish_date_desc",
		"format": "infinite",
	}, timePoint, timeDirection)

	if mergeParam != nil {
		param = maputil.Merge(param, mergeParam)
	}

	endpoint = strings.Trim(endpoint, "/")

	var result []model.Post
	var err error
	hasMore := true
	for hasMore {
		var moreList moreList[model.Post]
		err = c.req.GetUnmarshal(ApiURLPath("/users/%d/%s", userID, endpoint), param, &moreList)
		if err != nil {
			break
		}
		hasMore = moreList.HasMore
		result = append(result, moreList.List...)

		updatePublishTimeParam(param, timeDirection, moreList.moreMarker)
	}
	if len(result) != 0 {
		return result, nil
	}
	if err == nil {
		if !c.IsSubscribed(userID) {
			err = ErrorChangeAccount
		} else {
			err = errors.New("no posts found or subscription expired")
		}
	}
	return nil, err
}

func (c *OFAPI) GetCollectionsListUsers(listid string) ([]model.CollectionListUser, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.CollectionListUser
	for hasMore {
		var list []model.CollectionListUser
		err = c.req.GetUnmarshal("/lists/"+listid+"/users", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
		}, &list)
		if err != nil {
			break
		}
		hasMore = len(list) >= 50
		offset += len(list)
		result = append(result, list...)
	}
	if len(result) != 0 {
		return result, nil
	}
	return nil, err
}

func (c *OFAPI) GetCollections(filter ...CollectionFilter) ([]model.Collection, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.Collection
	for hasMore {
		var moreList moreList[model.Collection]
		err = c.req.GetUnmarshal("/lists", map[string]string{
			"offset":     strconv.Itoa(offset),
			"limit":      "50",
			"skip_users": "all",
			"format":     "infinite",
		}, &moreList)
		if err != nil {
			break
		}

		hasMore = moreList.HasMore
		offset += len(moreList.List)

		for _, collection := range moreList.List {
			if len(filter) != 0 && !filter[0](collection) {
				continue
			}
			result = append(result, collection)
		}
	}

	if len(result) != 0 {
		return result, nil
	}
	return nil, err
}

func (c *OFAPI) GetSubscriptions(subType SubscritionType, filter ...SubscribeFilter) ([]model.Subscription, error) {
	if subType == SubscritionTypeAll {
		var resultAll []model.Subscription
		subActivate, errActive := c.GetSubscriptions(SubscritionTypeActive, filter...)
		if errActive == nil {
			resultAll = append(resultAll, subActivate...)
		}
		subExpired, errExpired := c.GetSubscriptions(SubscritionTypeExpired, filter...)
		if errExpired == nil {
			resultAll = append(resultAll, subExpired...)
		}
		if len(resultAll) != 0 {
			return resultAll, nil
		}
		if errActive != nil {
			return nil, errActive
		}
		return nil, errExpired
	}

	var result []model.Subscription

	if len(filter) == 0 {
		filter = append(filter, SubscribeRestrictedFilter(false))
	}

	var err error
	hasMore := true
	offset := 0
	for hasMore {
		var moreList moreList[model.Subscription]
		err = c.req.GetUnmarshal("/subscriptions/subscribes", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
			"type":   string(subType),
			"format": "infinite",
		}, &moreList)
		if err != nil {
			break
		}
		hasMore = moreList.HasMore
		offset += len(moreList.List)

		for _, sub := range moreList.List {
			if !filter[0](sub) {
				continue
			}
			result = append(result, sub)
		}
	}

	if len(result) != 0 {
		result = slice.UniqueBy(result, func(sub model.Subscription) int64 {
			return sub.ID
		})
		return result, nil
	}
	return nil, err
}

func (c *OFAPI) GetUserByUsername(username string) (model.User, error) {
	return c.GetUser(username)
}

func (c *OFAPI) GetUserByID(userID int64) (model.User, error) {
	var um map[string]model.User
	err := c.req.GetUnmarshal(ApiURLPath("/users/list?x[]=%d", userID), nil, &um)
	if err != nil {
		return model.User{}, err
	}
	user, ok := um[strconv.FormatInt(userID, 10)]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

func (c *OFAPI) GetUser(userEndpoint string) (model.User, error) {
	var user model.User
	err := c.req.GetUnmarshal(ApiURLPath("/users/%s", userEndpoint), map[string]string{
		"limit": "50",
		"order": "publish_date_asc",
	}, &user)
	return user, err
}

func (c *OFAPI) GetFileInfo(url string) (common.HttpFileInfo, error) {
	return c.req.GetFileInfo(url)
}

func (c *OFAPI) IsSubscribed(userIdOrName any) bool {
	subs, err := c.GetSubscriptions(SubscritionTypeActive)
	if err != nil {
		return false
	}
	user := fmt.Sprintf("%v", userIdOrName)
	userName := user
	userID, err := strconv.ParseInt(user, 10, 64)
	if err == nil {
		userName = ""
	}
	return slice.ContainBy(subs, func(sub model.Subscription) bool {
		if sub.ID != 0 && sub.ID == userID {
			return true
		}
		if sub.Username != "" && sub.Username == userName {
			return true
		}
		return false
	})
}
