package baidu

import (
	"net/http"
)

const authUrl = "https://openapi.baidu.com/oauth/2.0/"
const apiHost = "https://pan.baidu.com/"

type AppConfig struct {
	ClientID     string
	ClientSecret string
	AppName      string
}

type RefreshToken struct {
	AppConfig
	RefreshKey  string
	AccessKey   string
	ExpiredTime int64
}

type Client struct {
	AppConfig
	RefreshToken string
	AccessToken  string
	ExpiredTime  int64

	http.Client
}

type File struct {
	FsID int    `json:"fs_id"`
	Dir  int    `json:"isdir"`
	Name string `json:"server_filename"`
	Path string `json:"path"`
}

type errorMsg struct {
	ErrorCode int    `json:"errno"`
	ErrorMsg  string `json:"errmsg"`
}
