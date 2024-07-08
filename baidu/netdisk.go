package baidu

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

func AuthorizeURL(clientID string) string {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", "oob")
	params.Add("scope", "basic,netdisk")
	return "https://openapi.baidu.com/oauth/2.0/authorize?" + params.Encode()
}

func NewByCode(appConfig AppConfig, code string) (*Client, error) {
	c := &Client{
		AppConfig: appConfig,
	}

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("redirect_uri", "oob")

	err := c.token(&values)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewByRefreshToken(appConfig AppConfig, refreshToken RefreshToken) (*Client, error) {
	c := &Client{
		AppConfig:    appConfig,
		RefreshToken: refreshToken.RefreshKey,
		AccessToken:  refreshToken.AccessKey,
		ExpiredTime:  refreshToken.ExpiredTime,
	}

	err := c.checkAccessToken()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) IsDirExist(dir string) (bool, error) {
	values := url.Values{}
	values.Set("method", "list")
	values.Set("access_token", c.AccessToken)
	values.Set("dir", dir)
	values.Set("start", "0")
	values.Set("limit", "1")

	body, err := c.do("GET", apiHost+"rest/2.0/xpan/file", &values)
	if err != nil {
		return false, err
	}

	r := new(errorMsg)
	err = json.Unmarshal(body, r)
	if err != nil {
		return false, err
	}

	if r.ErrorCode != 0 && r.ErrorCode != -9 {
		return false, errors.New(r.ErrorMsg)
	}

	if r.ErrorCode == -9 {
		return false, nil
	} else {
		return true, nil
	}
}

func (c *Client) GetAllFiles() ([]File, error) {
	values := url.Values{}
	values.Set("method", "listall")
	values.Set("access_token", c.AccessToken)
	values.Set("path", "/apps/"+c.AppName)
	values.Set("recursion", "1")

	cursor := 0
	hasMore := true
	var files []File
	for hasMore {
		values.Set("start", strconv.Itoa(cursor))

		body, err := c.do("GET", apiHost+"rest/2.0/xpan/multimedia", &values)
		if err != nil {
			return nil, err
		}

		type rsp struct {
			errorMsg
			Cursor  int    `json:"cursor"`
			HasMore int    `json:"has_more"`
			List    []File `json:"list"`
		}

		r := new(rsp)
		err = json.Unmarshal(body, r)
		if err != nil {
			return nil, err
		}

		cursor = r.Cursor
		hasMore = r.HasMore == 1

		files = append(files, r.List...)
	}

	return files, nil
}
