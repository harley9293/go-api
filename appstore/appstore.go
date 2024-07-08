package appstore

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const sandboxUrl = "https://sandbox.itunes.apple.com/"
const buyUrl = "https://buy.itunes.apple.com/"

func VerifyReceipt(secret, receipt string, sandbox bool) error {
	params := make(map[string]string)
	params["receipt-data"] = receipt
	params["password"] = secret
	_, err := do("verifyReceipt", &params, sandbox)
	return err
}

func do(path string, params *map[string]string, sandbox bool) ([]byte, error) {
	jsonStr, err := json.Marshal(params)
	if err != nil {
		return []byte{}, err
	}

	url := buyUrl
	if sandbox {
		url = sandboxUrl
	}

	request, err := http.NewRequest("POST", url+path, bytes.NewBuffer(jsonStr))
	if err != nil {
		return []byte{}, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request.WithContext(context.TODO()))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
