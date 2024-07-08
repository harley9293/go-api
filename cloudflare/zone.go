package cloudflare

import (
	"encoding/json"
	"errors"
	"github.com/harley9293/go-util/net"
	"net/http"
)

type ZoneInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ListZones() ([]ZoneInfo, error) {
	c.clearRequest()
	c.request.Method = "GET"
	c.request.Path = url + "/zones"
	code, body, err := net.ExecuteHttpRequest(c.request)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var ret struct {
		Result []ZoneInfo `json:"result"`
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Result, nil
}
