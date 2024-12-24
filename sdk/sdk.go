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

var tr *http.Transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func (c *Client) CreateAccount(request api.CreateAccountRequest) (account model.Account, err *api.Error) {
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

func (c *Client) Retrieve(token string) (model.Account, *api.Error) {
	body, _ := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: token})
	req, _ := http.NewRequest("POST", c.BaseUrl+"/sessions/retrieve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", c.AccessToken)
	client := &http.Client{Transport: tr}
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		var err api.Error
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &err)
		return model.Account{}, &err
	}
	var account model.Account
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &account)
	return account, nil
}

func (c *Client) DeleteAccount(namespace, id string) *api.Error {
	req, _ := http.NewRequest("DELETE", c.BaseUrl+"/account", nil)
	req.Header.Set("X-Access-Token", c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("namespace", namespace)
	req.Header.Set("id", id)
	req.Header.Set("idby", "id")
	client := &http.Client{Transport: tr}
	resp, _ := client.Do(req)
	if resp.StatusCode != 204 {
		var err api.Error
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &err)
		return &err
	}
	return nil
}

func (c *Client) ListAccounts(namespace string) []model.Pagination[model.Account] {
	req, _ := http.NewRequest("GET", c.BaseUrl+"/accounts", nil)
	req.Header.Set("X-Access-Token", c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("namespace", namespace)
	client := &http.Client{Transport: tr}
	resp, _ := client.Do(req)
	var accounts []model.Pagination[model.Account]
	if resp.StatusCode != 200 {
		return accounts
	}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &accounts)
	return accounts
}
