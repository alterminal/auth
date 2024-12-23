package repo

import (
	"github.com/alterminal/auth/model"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.AccountEmail{})
	db.AutoMigrate(&model.AccountTable{})
	db.AutoMigrate(&model.AccountPhone{})
}
