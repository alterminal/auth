package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/macaoservices/auth/model"
	"github.com/macaoservices/auth/repo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(viper.GetString("jwt.secret"))
}

func Run(port string, db *gorm.DB) {
	router := gin.Default()
	router.Use(GetAccountMiddleware(db))
	{
		group := router.Group("accounts")
		group.GET("me", GetAccountMe(db))
		group.POST("sessions", CreateAccountSessions(db))
	}
	{
		group := router.Group("organizations")
		group.POST("", CreateOrganization(db))
	}
	router.Run(":" + port)
}

func GetAccountMiddleware(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		splitToken := strings.Split(authorization, " ")
		if len(splitToken) != 2 {
			return
		}
		tokenString := splitToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			return
		}
		id := token.Claims.(jwt.MapClaims)["id"]
		account, err := repo.FindAccount(db, repo.ByID(id.(string)))
		if err != nil {
			return
		}
		ctx.Set("account", account)
	}
}

func GetAccount(ctx *gin.Context) (*model.Account, bool) {
	account, ok := ctx.Get("account")
	if !ok {
		return nil, false
	}
	acc, ok := account.(*model.Account)
	if !ok {
		return nil, false
	}
	return acc, true
}

func GetAccountMe(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println(GetAccount(c))
	}
}

func CreateAccountSessions(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		req, ok := ShouldBindJSON[CreateSessionRequest](c)
		if !ok {
			return
		}
		account, err := repo.FindAccount(db, repo.ByAccount(req.Account))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if !account.CheckPassword(req.Password) {
			c.JSON(400, gin.H{"error": "invalid password"})
			return
		}
		claims := jwt.MapClaims{
			"id":    account.ID,
			"email": account.Email,
			"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)
		c.JSON(201, gin.H{"token": tokenString})
	}
}

func CreateOrganization(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		req, ok := ShouldBindJSON[UpsertOrganizationRequest](c)
		if !ok {
			return
		}
		org := model.Organization{
			Name: req.Name,
		}
		db.Create(&org)
	}
}
