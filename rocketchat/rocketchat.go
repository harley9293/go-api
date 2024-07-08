package rocketchat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Room struct {
	Id    string `json:"_id"`
	Name  string `json:"name"`
	FName string `json:"fname"`
}

type RC struct {
	url       string
	userId    string
	authToken string
	login     bool
	roomList  []Room
}

func New(url string) *RC {
	return &RC{url: url + "/api/v1/", login: false}
}

func (rc *RC) Login(user, password string) error {
	params := make(map[string]string)
	params["user"] = user
	params["password"] = password
	rsp, err := rc.do("POST", "login", &params)
	if err != nil {
		return err
	}

	type LoginRspData struct {
		UserId    string `json:"userId"`
		AuthToken string `json:"authToken"`
	}

	type LoginRsp struct {
		Status string       `json:"status"`
		Data   LoginRspData `json:"data"`
	}

	loginRsp := new(LoginRsp)
	err = json.Unmarshal(rsp, loginRsp)
	if err != nil {
		return err
	}

	rc.userId = loginRsp.Data.UserId
	rc.authToken = loginRsp.Data.AuthToken
	rc.login = true

	err = rc.initRoomList()
	if err != nil {
		return err
	}

	return nil
}

func (rc *RC) PostMessage(text string, roomName string) error {
	params := make(map[string]string)
	params["roomId"] = rc.findRoomIdByFName(roomName)
	params["text"] = text
	rsp, err := rc.do("POST", "chat.postMessage", &params)
	if err != nil {
		return err
	}

	type PostMessageRsp struct {
		Success bool `json:"success"`
	}

	postMessageRsp := new(PostMessageRsp)
	err = json.Unmarshal(rsp, postMessageRsp)
	if err != nil {
		return err
	}

	if postMessageRsp.Success == false {
		return errors.New(string(rsp))
	}

	return nil
}

func (rc *RC) findRoomIdByFName(fName string) string {
	for _, room := range rc.roomList {
		if room.FName == fName {
			return room.Id
		}
	}

	return ""
}

func (rc *RC) initRoomList() error {
	params := make(map[string]string)
	rsp, err := rc.do("GET", "rooms.get", &params)
	if err != nil {
		return err
	}

	type RoomListRsp struct {
		Success bool   `json:"success"`
		Update  []Room `json:"update"`
	}

	roomListRsp := new(RoomListRsp)
	err = json.Unmarshal(rsp, roomListRsp)
	if err != nil {
		return err
	}

	if roomListRsp.Success == false {
		return errors.New(string(rsp))
	}

	rc.roomList = roomListRsp.Update

	return nil
}

func (rc *RC) do(method string, path string, params *map[string]string) ([]byte, error) {
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return []byte{}, err
	}

	request, err := http.NewRequest(method, rc.url+path, bytes.NewBuffer(jsonStr))
	if err != nil {
		return []byte{}, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	if rc.login {
		request.Header.Set("X-Auth-Token", rc.authToken)
		request.Header.Set("X-User-Id", rc.userId)
	}
	client := http.Client{}
	resp, err := client.Do(request.WithContext(context.TODO()))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
