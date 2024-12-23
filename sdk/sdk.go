package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alterminal/auth/api"
	"github.com/alterminal/auth/model"
)

type Client struct {
	BaseUrl     string
	AccessToken string
}

func (c *Client) CreateAccount(request api.CreateAccountRequest) (account model.Account, err *api.Error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", c.BaseUrl+"/account", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", c.AccessToken)
	resp, _ := client.Do(req)
	if resp.StatusCode != 201 {
		var err api.Error
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &err)
		return account, &err
	}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &account)
	return account, nil
}

func (c *Client) CreateSession(request api.CreateSessionRequest) (token string, err *api.Error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", c.BaseUrl+"/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", c.AccessToken)
	resp, _ := client.Do(req)
	if resp.StatusCode != 201 {
		var err api.Error
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &err)
		return token, &err
	}
	var respToken struct {
		Token string `json:"token"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respToken)
	return respToken.Token, nil
}
