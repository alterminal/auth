package api

import (
	"github.com/alterminal/common/mid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	router := gin.Default()
	router.Use(mid.AccessControllAllowfunc(mid.AccessControllAllowConfig{
		Origin:  "*",
		Headers: "*",
		Methods: "*",
	}))
	accountApi := AccountApi{db: db}
	accountApi.BindRouter(router)
	router.Run(":8080")
}

type AccountApi struct {
	db *gorm.DB
}

func (api *AccountApi) BindRouter(router *gin.Engine) {
	router.POST("/account", api.CreateAccount)
}

func (api *AccountApi) CreateAccount(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
