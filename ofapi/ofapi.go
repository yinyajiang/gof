package ofapi

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi/model"
)

type OFApi struct {
	cfg OFApiConfig
}

type OFApiConfig struct {
	AuthInfo gof.AuthInfo
	Rules    gof.Rules
}

func NewOFAPI(config OFApiConfig) *OFApi {
	return &OFApi{cfg: config}
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

type UserIdentifier struct {
	ID       int64
	Username string
}

func (c *OFApi) GetPaidPosts() ([]model.PaidPost, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.PaidPost

	for hasMore {
		var moreList MoreList[model.PaidPost]
		err = OFApiAuthGetUnmashel("/posts/paid", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
			"order":  "publish_date_desc",
			"format": "infinite",
		}, c.cfg.AuthInfo, c.cfg.Rules, &moreList)
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
	err = OFApiAuthGetUnmashel(ApiURLPath("/posts/%s", postURLInfo.PostID), map[string]string{
		"skip_users": "all",
	}, c.cfg.AuthInfo, c.cfg.Rules, &post)
	return post, err
}

func (c *OFApi) GetCollectionsListUsers(listid string) ([]model.CollectionListUser, error) {
	var err error
	hasMore := true
	offset := 0
	var result []model.CollectionListUser
	for hasMore {
		var moreList MoreList[model.CollectionListUser]
		err = OFApiAuthGetUnmashel("/lists/"+listid+"/users", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
		}, c.cfg.AuthInfo, c.cfg.Rules, &moreList)
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
		var moreList MoreList[model.Collection]
		err = OFApiAuthGetUnmashel("/lists", map[string]string{
			"offset":     strconv.Itoa(offset),
			"limit":      "50",
			"skip_users": "all",
			"format":     "infinite",
		}, c.cfg.AuthInfo, c.cfg.Rules, &moreList)
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
		var moreList MoreList[model.Subscription]
		err = OFApiAuthGetUnmashel("/subscriptions/subscribes", map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  "50",
			"type":   string(subType),
			"format": "infinite",
		}, c.cfg.AuthInfo, c.cfg.Rules, &moreList)
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
	data, err := OFApiAuthGet(ApiURLPath("/users/%s", userEndpoint), map[string]string{
		"limit": "50",
		"order": "publish_date_asc",
	}, c.cfg.AuthInfo, c.cfg.Rules)
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

type MoreList[T any] struct {
	HasMore bool `json:"hasMore"`
	List    []T  `json:"list"`
}
