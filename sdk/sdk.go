package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/alterminal/auth/api"
	"github.com/alterminal/auth/model"
)

type Client struct {
	BaseUrl     string
	AccessToken string
}

type AccountOption func(url.Values) url.Values

func WithId(id string) AccountOption {
	return func(v url.Values) url.Values {
		v.Set("idby", "id")
		v.Add("id", id)
		return v
	}
}

func WithEmail(email string) AccountOption {
	return func(v url.Values) url.Values {
		v.Set("idby", "email")
		v.Add("email", email)
		return v
	}
}

func WithPhone(phoneRegion, phoneNumber string) AccountOption {
	return func(v url.Values) url.Values {
		v.Set("idby", "phone")
		v.Add("phoneRegion", phoneRegion)
		v.Add("phoneNumber", phoneNumber)
		return v
	}
}

var tr *http.Transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

func (c *Client) CreateAccount(request api.CreateAccountRequest) (account model.Account, err *api.Error) {
	client := &http.Client{Transport: tr}
	body, _ := json.Marshal(request)
	u, _ := url.ParseRequestURI(c.BaseUrl + "/account")
	req := c.buildRequest("POST", u, bytes.NewReader(body))
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
	u, _ := url.ParseRequestURI(c.BaseUrl + "/sessions")
	req := c.buildRequest("POST", u, bytes.NewReader(body))
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
	u, _ := url.ParseRequestURI(c.BaseUrl + "/sessions/retrieve")
	req := c.buildRequest("POST", u, bytes.NewReader(body))
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

func (c *Client) DeleteAccount(namespace string, option AccountOption) *api.Error {
	paramsValue := url.Values{}
	paramsValue.Add("namespace", namespace)
	paramsValue = option(paramsValue)
	u, _ := url.ParseRequestURI(c.BaseUrl + "/account")
	u.RawQuery = paramsValue.Encode()
	req := c.buildRequest("DELETE", u, nil)
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

func (c *Client) ListAccounts(namespace string) model.Pagination[*model.Account] {
	paramsValue := url.Values{}
	paramsValue.Add("namespace", namespace)
	u, _ := url.ParseRequestURI(c.BaseUrl + "/accounts")
	u.RawQuery = paramsValue.Encode()
	req := c.buildRequest("GET", u, nil)
	client := &http.Client{Transport: tr}
	resp, _ := client.Do(req)
	var accounts model.Pagination[*model.Account]
	if resp.StatusCode != 200 {
		return accounts
	}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &accounts)
	return accounts
}

func (c *Client) GetAccount(namespace string, option AccountOption) (model.Account, *api.Error) {
	paramsValue := url.Values{}
	paramsValue = option(paramsValue)
	paramsValue.Add("namespace", namespace)
	u, _ := url.ParseRequestURI(c.BaseUrl + "/account")
	u.RawQuery = paramsValue.Encode()
	req := c.buildRequest("GET", u, nil)
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

func (c *Client) SetPassword(namespace string, option AccountOption) *api.Error {
	paramsValue := url.Values{}
	paramsValue = option(paramsValue)
	paramsValue.Add("namespace", namespace)
	u, _ := url.ParseRequestURI(c.BaseUrl + "/account/password")
	u.RawQuery = paramsValue.Encode()
	req := c.buildRequest("PUT", u, nil)
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

func (c *Client) buildRequest(method string, u *url.URL, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("%v", u), body)
	req.Header.Set("X-Access-Token", c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	return req
}
