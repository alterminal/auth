package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestModel(t *testing.T) {
	var account Account
	account.ID = "1234567890123456789"
	account.Account = "account"
	account.Email = "email"
	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now
	t.Log(now.Unix())
	j, _ := json.Marshal(account)
	t.Log(string(j))
	var account2 Account
	json.Unmarshal(j, &account2)
	t.Log(account, account2)
}
