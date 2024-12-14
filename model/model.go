package model

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Account{}, &Organization{}, &Role{}, &Member{}, &Permission{})
	var org Organization
	org.ID = viper.GetString("system.organization.id")
	org.Name = viper.GetString("system.organization.name")
	err := db.FirstOrCreate(&org).Error
	if err != nil {
		panic(err)
	}
	role := Role{
		ID:             "0000000000000000000",
		OrganizationID: org.ID,
		Title:          "root",
	}
	db.FirstOrCreate(&role)
	accounts, ok := viper.Get("system.accounts").([]interface{})
	if ok {
		for _, account := range accounts {
			a := account.(map[string]interface{})
			var acc Account
			acc.ID = a["id"].(string)
			acc.Account = a["account"].(string)
			acc.Email = a["email"].(string)
			acc.SetPassword(a["password"].(string))
			db.FirstOrCreate(&acc)
			member := Member{
				AccountID: acc.ID,
				RoleID:    role.ID,
			}
			db.FirstOrCreate(&member)
		}
	}
}
