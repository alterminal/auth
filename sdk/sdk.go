package sdk

import "github.com/alterminal/auth/model"

type Client struct {
	BaseUrl     string
	AccessToken string
	Namespace   *string
}

func (c *Client) CreateAccount(account *model.Account) (*model.Account, error) {
	return nil, nil
}
