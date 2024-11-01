package ofapi

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi/model"
)

type OFApi struct {
	req *Req
}

type OFApiConfig struct {
	AuthInfo gof.AuthInfo
	Rules    gof.Rules
}

func NewOFAPI(config OFApiConfig) *OFApi {
	common.PanicAuthInfo(config.AuthInfo)
	return &OFApi{
		req: &Req{
			authInfo: config.AuthInfo,
			rules:    config.Rules,
		},
	}
}

func (c *OFApi) Req() *Req {
	return c.req
}

func (c *OFApi) CheckAuth() error {
	me, err := c.GetMe()
	if err != nil {
		return err
	}
	if me.Username == "" && me.Name == "" {
		return errors.New("auth failed")
	}
	return nil
}

func (c *OFApi) GetMe() (model.User, error) {
	return c.GetUser("me")
}

func (c *OFApi) GetUserHightlights(userID int64, withStories ...bool) ([]model.Highlight, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.Highlight
	for hasMore {
		var moreList moreList[model.Highlight]
		err = c.req.GetUnmashel(ApiURLPath("/users/%d/stories/highlights", userID), map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "5",
		}, &moreList)
		if err != nil {
			break
		}

		hasMore = moreList.HasMore
		offset += len(moreList.List)

		if len(withStories) != 0 && withStories[0] {
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

func (c *OFApi) GetHighlight(highlightID int64) (model.Highlight, error) {
	var highlight model.Highlight
	err := c.req.GetUnmashel(ApiURLPath("/stories/highlights/%d", highlightID), nil, &highlight)
	return highlight, err
}

func (c *OFApi) GetUserStories(userID int64) ([]model.Story, error) {
	var stories []model.Story
	err := c.req.GetUnmashel(ApiURLPath("/users/%d/stories", userID), map[string]string{
		"limit": "50",
		"order": "publish_date_desc",
	}, &stories)
	return stories, err
}

func (c *OFApi) GetPaidPosts() ([]model.Post, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.Post

	for hasMore {
		var moreList moreList[model.Post]
		err = c.req.GetUnmashel("/posts/paid", map[string]string{
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

func (c *OFApi) GetPost(postURL string) (model.Post, error) {
	postURLInfo, err := common.ParseSinglePostURL(postURL)
	if err != nil {
		return model.Post{}, err
	}
	var post model.Post
	err = c.req.GetUnmashel(ApiURLPath("/posts/%s", postURLInfo.PostID), map[string]string{
		"skip_users": "all",
	}, &post)
	return post, err
}

func (c *OFApi) GetUserPosts(userID int64) ([]model.Post, error) {
	return c.GetUserPostsByTime(userID, time.Now(), TimeDirectionBefore)
}

func (c *OFApi) GetUserMedias(userID int64) ([]model.Post, error) {
	return c.GetUserMediasByTime(userID, time.Now(), TimeDirectionBefore)
}

func (c *OFApi) GetUserMediasByTime(userID int64, timePoint time.Time, timeDirection TimeDirection) ([]model.Post, error) {
	param := initPublishTimeParam(map[string]string{
		"limit":      "50",
		"order":      "publish_date_desc",
		"format":     "infinite",
		"skip_users": "all",
	}, timePoint, timeDirection)

	var result []model.Post
	var err error
	hasMore := true
	for hasMore {
		var moreList moreList[model.Post]
		err = c.req.GetUnmashel(ApiURLPath("/users/%d/posts/medias", userID),
			param, &moreList)
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
	return nil, err
}

func (c *OFApi) GetUserPostsByTime(userID int64, timePoint time.Time, timeDirection TimeDirection) ([]model.Post, error) {
	param := initPublishTimeParam(map[string]string{
		"limit":  "50",
		"order":  "publish_date_desc",
		"format": "infinite",
	}, timePoint, timeDirection)

	var result []model.Post
	var err error
	hasMore := true
	for hasMore {
		var moreList moreList[model.Post]
		err = c.req.GetUnmashel(ApiURLPath("/users/%d/posts", userID), param, &moreList)
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
	return nil, err
}

func (c *OFApi) GetCollectionsListUsers(listid string) ([]model.CollectionListUser, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.CollectionListUser
	for hasMore {
		var moreList moreList[model.CollectionListUser]
		err = c.req.GetUnmashel("/lists/"+listid+"/users", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
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

func (c *OFApi) GetCollections(filter ...CollectionFilter) ([]model.Collection, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.Collection
	for hasMore {
		var moreList moreList[model.Collection]
		err = c.req.GetUnmashel("/lists", map[string]string{
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

func (c *OFApi) GetSubscriptions(subType SubscritionType, filter ...SubscribeFilter) ([]model.Subscription, error) {
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
		err = c.req.GetUnmashel("/subscriptions/subscribes", map[string]string{
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

func (c *OFApi) GetUserByUsername(username string) (model.User, error) {
	return c.GetUser(username)
}

func (c *OFApi) GetUser(userEndpoint string) (model.User, error) {
	data, err := c.req.Get(ApiURLPath("/users/%s", userEndpoint), map[string]string{
		"limit": "50",
		"order": "publish_date_asc",
	})
	if err != nil {
		return model.User{}, err
	}
	var userInfo model.User
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		return model.User{}, err
	}
	return userInfo, nil
}
