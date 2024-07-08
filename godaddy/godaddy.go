package godaddy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const url = "https://api.godaddy.com"

type Record struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Type string `json:"type"`
	Ttl  int    `json:"ttl"`
}

type Client struct {
	key    string
	secret string
}

func New(key, secret string) *Client {
	return &Client{
		key:    key,
		secret: secret,
	}
}

func (c *Client) GetDomainRecords(domain, ty, name string) ([]Record, error) {
	url := "/v1/domains/" + domain + "/records"
	if ty != "" {
		url += "/" + ty
		if name != "" {
			url += "/" + name
		}
	}

	body, err := c.do("GET", url, nil)
	if err != nil {
		return []Record{}, err
	}

	var ret []Record
	json.Unmarshal(body, &ret)
	return ret, nil
}

func (c *Client) UpdateDomainRecords(domain, ty, name, data string) error {
	url := "/v1/domains/" + domain + "/records/" + ty + "/" + name

	record := Record{
		Data: data,
		Name: name,
		Type: ty,
		Ttl:  600,
	}

	_, err := c.do("PUT", url, []interface{}{record})
	return err
}

func (c *Client) do(method string, path string, params []interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url+path, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "sso-key "+c.key+":"+c.secret)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Connection", "close")
	client := http.Client{}
	resp, err := client.Do(request.WithContext(context.TODO()))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(resp.StatusCode))
	}
	return body, nil
}
