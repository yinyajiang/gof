package ofapi

import (
	"encoding/json"

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
	_, err := c.GetMeUserInfo()
	return err
}

func (c *OFApi) GetMeUserInfo() (model.UserInfo, error) {
	return c.GetUserInfo("/users/me")
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
