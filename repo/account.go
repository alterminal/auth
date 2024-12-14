package repo

import (
	"errors"

	"github.com/macaoservices/auth/fault"
	"github.com/macaoservices/auth/model"
	"gorm.io/gorm"
)

type Option func(*gorm.DB) *gorm.DB

func FindOne[T any](tx *gorm.DB, options ...Option) (*T, error) {
	var t T
	ctx := tx
	for _, option := range options {
		ctx = option(ctx)
	}
	err := tx.First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fault.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func ByID(id string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id)
	}
}

func ByEmail(email string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("email = ?", email)
	}
}

func ByAccount(account string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("account = ?", account)
	}
}

func FindAccount(tx *gorm.DB, options ...Option) (*model.Account, error) {
	return FindOne[model.Account](tx, options...)
}
