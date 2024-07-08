package cloudflare

import (
	"github.com/harley9293/go-util/net"
)

const url = "https://api.cloudflare.com/client/v4"

type Client struct {
	request net.HttpRequestParams
}

func NewClient(bearerAuth string) *Client {
	return &Client{
		request: net.HttpRequestParams{
			Headers: map[string]string{
				"Authorization": "Bearer " + bearerAuth,
				"Content-Type":  "application/json;charset=UTF-8",
			},
		},
	}
}

func (c *Client) clearRequest() {
	c.request.Method = ""
	c.request.Path = ""
	c.request.Params = nil
}
