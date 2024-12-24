package model

import (
	"time"

	"github.com/alterminal/auth/utils"
	"github.com/alterminal/common/utils/pwd"
	"gorm.io/gorm"
)

type Account struct {
	Namespace string    `gorm:"type:varchar(64);primaryKey" json:"namespace"`
	ID        string    `gorm:"type:char(19);primaryKey" json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	// virtual field
	Account     *string `json:"account" gorm:"-"`
	Email       *string `json:"email" gorm:"-"`
	PhoneRegion *string `json:"phoneRegion" gorm:"-"`
	PhoneNumber *string `json:"phoneNumber" gorm:"-"`
	// hidden field
	Password string `json:"-" gorm:"type:char(128)"`
	Salt     string `json:"-" gorm:"type:char(16)"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.Account != nil {
		err := tx.Create(&AccountTable{Namespace: a.Namespace, Account: *a.Account, ID: a.ID}).Error

		if err != nil {
			return err
		}
	}
	if a.Email != nil {
		err := tx.Create(&AccountEmail{Namespace: a.Namespace, Email: *a.Email, ID: a.ID}).Error
		if err != nil {
			return err
		}
	}
	if a.PhoneRegion != nil && a.PhoneNumber != nil {
		err := tx.Create(&AccountPhone{Namespace: a.Namespace, PhoneRegion: *a.PhoneRegion, PhoneNumber: *a.PhoneNumber, ID: a.ID}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Account) CheckPassword(password string) bool {
	return pwd.CheckPassword(a.Password, password, a.Salt)
}

func (a *Account) SetPassword(password string) {
	a.Password, a.Salt = utils.HashWithSalt(password)
}

func (a *Account) BeforeDelete(tx *gorm.DB) error {
	tx.Where("namespace = ? AND id = ?", a.Namespace, a.ID).Delete(&AccountEmail{})
	tx.Where("namespace = ? AND id = ?", a.Namespace, a.ID).Delete(&AccountTable{})
	tx.Where("namespace = ? AND id = ?", a.Namespace, a.ID).Delete(&AccountPhone{})
	return nil
}

func (a *Account) AfterFind(tx *gorm.DB) error {
	var email AccountEmail
	if err := tx.First(&email, "namespace = ? AND id = ?", a.Namespace, a.ID).Error; err == nil {
		a.Email = &email.Email
	}
	var acc AccountTable
	if err := tx.First(&acc, "namespace = ? AND id = ?", a.Namespace, a.ID).Error; err == nil {
		a.Account = &acc.Account
	}
	var phone AccountPhone
	if err := tx.First(&phone, "namespace = ? AND id = ?", a.Namespace, a.ID).Error; err == nil {
		a.PhoneRegion = &phone.PhoneRegion
		a.PhoneNumber = &phone.PhoneNumber
	}
	return nil
}

type AccountEmail struct {
	Namespace string `gorm:"type:varchar(64);primaryKey" json:"namespace"`
	Email     string `gorm:"type:varchar(128);primaryKey" json:"email"`
	ID        string `gorm:"type:char(19);uniqueIndex" json:"id"`
}

type AccountTable struct {
	Namespace string `gorm:"type:varchar(64);primaryKey" json:"namespace"`
	Account   string `gorm:"type:varchar(128);primaryKey" json:"account"`
	ID        string `gorm:"type:char(19);uniqueIndex" json:"id"`
}

type AccountPhone struct {
	Namespace   string `gorm:"type:varchar(64);primaryKey" json:"namespace"`
	PhoneRegion string `gorm:"type:char(8);primaryKey" json:"phoneRegion"`
	PhoneNumber string `gorm:"type:char(11);primaryKey" json:"phoneNumber"`
	ID          string `gorm:"type:char(19);uniqueIndex" json:"id"`
}
