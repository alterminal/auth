package model

import (
	"encoding/json"
	"time"

	"github.com/macaoservices/auth/utils"
	"gorm.io/gorm"
)

type Account struct {
	ID          string         `json:"id" gorm:"type:char(19);primaryKey"`
	Account     string         `json:"account" gorm:"type:varchar(64);uniqueIndex"`
	PhoneNumber string         `json:"phoneNumber" gorm:"type:varchar(11);uniqueIndex"`
	Email       string         `json:"email" gorm:"type:varchar(64);uniqueIndex"`
	Password    string         `json:"-" gorm:"type:char(128)"`
	Salt        string         `json:"-" gorm:"type:char(16)"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (a *Account) CheckPassword(password string) bool {
	return utils.CheckPassword(a.Password, password, a.Salt)
}

func (a *Account) SetPassword(password string) {
	hashed, salt := utils.HashWithSalt(password)
	a.Password = hashed
	a.Salt = salt
}

func (a *Account) UnmarshalJSON(data []byte) error {
	jsonMap := make(map[string]any)
	json.Unmarshal(data, &jsonMap)
	if v, ok := jsonMap["id"]; ok {
		a.ID = v.(string)
	}
	if v, ok := jsonMap["account"]; ok {
		a.Account = v.(string)
	}
	if v, ok := jsonMap["email"]; ok {
		a.Email = v.(string)
	}
	if v, ok := jsonMap["createdAt"]; ok {
		timestamp, ok := v.(float64)
		if ok {
			a.CreatedAt = time.Unix(int64(timestamp), 0)
		}
	}
	if v, ok := jsonMap["updatedAt"]; ok {
		timestamp, ok := v.(float64)
		if ok {
			a.UpdatedAt = time.Unix(int64(timestamp), 0)
		}
	}
	return nil
}

func (a Account) MarshalJSON() ([]byte, error) {
	jsonMap := make(map[string]any)
	jsonMap["id"] = a.ID
	jsonMap["account"] = a.Account
	jsonMap["email"] = a.Email
	jsonMap["createdAt"] = a.CreatedAt.Unix()
	jsonMap["updatedAt"] = a.UpdatedAt.Unix()
	return json.Marshal(jsonMap)
}
