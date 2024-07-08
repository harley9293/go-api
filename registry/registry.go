package registry

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct {
	url      string
	client   *http.Client
	username string
	password string
}

func New(url, username, password string) (*Client, error) {
	r := &Client{
		url:      url,
		client:   &http.Client{},
		username: username,
		password: password,
	}

	_, err := r.do("GET", "/v2/")
	return r, err
}

func (c *Client) Repositories() ([]string, error) {
	type data struct {
		Names []string `json:"repositories"`
	}

	body, err := c.do("GET", "/v2/_catalog")
	d := &data{}
	err = json.Unmarshal(body, d)
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, name := range d.Names {
		repos = append(repos, name)
	}

	return repos, err
}

func (c *Client) do(method string, path string) ([]byte, error) {
	request, err := http.NewRequest(method, c.url+path, nil)
	if err != nil {
		return nil, err
	}

	if c.username != "" || c.password != "" {
		request.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.client.Do(request.WithContext(context.TODO()))
	if err != nil {
		return nil, err
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
