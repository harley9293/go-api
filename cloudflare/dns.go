package cloudflare

import (
	"encoding/json"
	"errors"
	"github.com/harley9293/go-util/net"
	"net/http"
)

type DNSInfo struct {
	ID      string `json:"id"`
	TTL     int    `json:"ttl"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

func (c *Client) ListDNSRecords(zoneID string) ([]DNSInfo, error) {
	c.clearRequest()
	c.request.Method = "GET"
	c.request.Path = url + "/zones/" + zoneID + "/dns_records"
	code, body, err := net.ExecuteHttpRequest(c.request)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var ret struct {
		Result []DNSInfo `json:"result"`
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, err
	}

	return ret.Result, nil
}

func (c *Client) UpdateDNSRecord(zoneID string, dnsInfo DNSInfo) error {
	c.clearRequest()
	c.request.Method = "PATCH"
	c.request.Path = url + "/zones/" + zoneID + "/dns_records/" + dnsInfo.ID
	c.request.Params = map[string]interface{}{
		"type":    dnsInfo.Type,
		"name":    dnsInfo.Name,
		"content": dnsInfo.Content,
		"ttl":     dnsInfo.TTL,
	}
	code, body, err := net.ExecuteHttpRequest(c.request)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}
