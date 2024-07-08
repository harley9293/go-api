package baidu

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) token(values *url.Values) error {
	values.Add("client_id", c.ClientID)
	values.Add("client_secret", c.ClientSecret)

	body, err := c.do("GET", authUrl+"token", values)
	if err != nil {
		return err
	}

	type rsp struct {
		ExpiredTime  int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}

	r := new(rsp)
	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	c.RefreshToken = r.RefreshToken
	c.AccessToken = r.AccessToken
	c.ExpiredTime = time.Now().Unix() + r.ExpiredTime

	return nil
}

func (c *Client) checkAccessToken() error {
	if c.ExpiredTime <= time.Now().Unix() {
		values := url.Values{}
		values.Add("grant_type", "refresh_token")
		values.Add("refresh_token", c.RefreshToken)
		return c.token(&values)
	}
	return nil
}

func (c *Client) do(method string, url string, values *url.Values) ([]byte, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = values.Encode()

	resp, err := c.Do(request.WithContext(context.TODO()))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	if description, ok := m["error_description"]; ok {
		return nil, errors.New(description.(string))
	}

	return body, nil
}
