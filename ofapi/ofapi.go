package ofapi

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/yinyajiang/gof"
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
	me, err := c.GetMeUserInfo()
	if err != nil {
		return err
	}
	if me.Username == "" && me.Name == "" {
		return errors.New("auth failed")
	}
	return nil
}

func (c *OFApi) GetMeUserInfo() (model.UserInfo, error) {
	return c.GetUserInfo("/users/me")
}

func (c *OFApi) GetSubscriptions(subType SubscritionType, filters ...SubscribeFilter) ([]model.Subscription, error) {
	if len(filters) == 0 {
		filters = append(filters, SubscribeRestrictedFilter(false))
	}

	var subscriptions []model.Subscription

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
			if !filters[0](sub) {
				continue
			}
			subscriptions = append(subscriptions, sub)
		}
	}
	if len(subscriptions) != 0 {
		return subscriptions, nil
	}
	return nil, err
}

func (c *OFApi) GetUserInfo(endpoint string) (model.UserInfo, error) {
	data, err := OFApiAuthGet(endpoint, map[string]string{
		"limit": "50",
		"order": "publish_date_asc",
	}, c.cfg.AuthInfo, c.cfg.Rules)
	if err != nil {
		return model.UserInfo{}, err
	}
	var userInfo model.UserInfo
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		return model.UserInfo{}, err
	}
	return userInfo, nil
}

type MoreList[T any] struct {
	HasMore bool `json:"hasMore"`
	List    []T  `json:"list"`
}
