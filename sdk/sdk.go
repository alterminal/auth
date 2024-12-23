package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alterminal/auth/api"
	"github.com/alterminal/auth/model"
)

type Client struct {
	BaseUrl     string
	AccessToken string
}

func (c *Client) CreateAccount(request api.CreateAccountRequest) (account model.Account, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", c.BaseUrl+"/account", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Token", c.AccessToken)
	resp, _ := client.Do(req)
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &account)
	fmt.Println(resp)
	return account, nil
}
