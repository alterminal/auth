package api

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alterminal/auth/model"
	"github.com/alterminal/common/mid"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	port := "8080"
	if p := viper.GetString("port"); p != "" {
		port = p
	}
	accessToken := viper.GetString("X-Access-Token")
	router := gin.Default()
	router.Use(mid.AccessControllAllowfunc(mid.AccessControllAllowConfig{
		Origin:  "*",
		Headers: "*",
		Methods: "*",
	}))
	accountApi := AccountApi{db: db, AccessToken: accessToken}
	accountApi.BindRouter(router)
	router.Run(":" + port)
}

func AuthMiddleware(accessToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Access-Token")
		if token != accessToken {
			c.JSON(401, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
		}
	}
}

type AccountApi struct {
	db          *gorm.DB
	AccessToken string
}

func (api *AccountApi) BindRouter(router *gin.Engine) {
	router.GET("account", AuthMiddleware(api.AccessToken), api.GetAccount)
	router.POST("account", AuthMiddleware(api.AccessToken), api.CreateAccount)
	router.PUT("account/password", AuthMiddleware(api.AccessToken), api.UpdatePassword)
	router.GET("accounts", AuthMiddleware(api.AccessToken), api.ListAccount)
	router.DELETE("account", AuthMiddleware(api.AccessToken), api.DeleteAccount)
	router.POST("sessions", api.CreateSession)
	router.POST("sessions/retrieve", api.RetrieveSession)
}

func (api *AccountApi) GetAccount(ctx *gin.Context) {
	account, ok := GetAccount(api.db, ctx)
	if !ok {
		ctx.JSON(404, Error{
			Message:    "account not found",
			StatusCode: 404,
			Code:       "404001",
		})
		return
	}
	ctx.JSON(200, account)
}

func (api *AccountApi) RetrieveSession(ctx *gin.Context) {
	req, ok := ShouldBindJSON[struct {
		Token string `json:"token"`
	}](ctx)
	if !ok {
		return
	}
	claims, err := ParseJWT(req.Token)
	if err != nil {
		ctx.JSON(401, Error{
			Message:    "token invalid",
			StatusCode: 401,
			Code:       "401001",
		})
		return
	}
	var account model.Account
	err = api.db.First(&account, "namespace = ? AND id = ?", claims["namespace"], claims["id"]).Error
	if err != nil {
		ctx.JSON(404, Error{
			Message:    "account not found",
			StatusCode: 404,
			Code:       "404001",
		})
	}
	ctx.JSON(200, account)
}

func (api *AccountApi) DeleteAccount(ctx *gin.Context) {
	account, ok := GetAccount(api.db, ctx)
	if !ok {
		ctx.JSON(404, Error{
			Message:    "account not found",
			StatusCode: 404,
			Code:       "404001",
		})
		return
	}
	err := api.db.Delete(account).Error
	if err != nil {
		ctx.JSON(500, Error{
			Message:    err.Error(),
			StatusCode: 500,
			Code:       "500001",
		})
		return
	}
	ctx.Status(204)
}

func (api *AccountApi) ListAccount(ctx *gin.Context) {
	namespace := ctx.Query("namespace")
	var page int64 = 0
	var limit int64 = 10
	if limitString := ctx.Query("limit"); limitString != "" {
		limit, _ = strconv.ParseInt(limitString, 10, 64)
	}
	if pageString := ctx.Query("page"); pageString != "" {
		page, _ = strconv.ParseInt(pageString, 10, 64)
	}
	accountList, err := model.ListByOption[model.Account](api.db, int(limit), int(page), model.WithNamespace(namespace))
	if err != nil {
		ctx.JSON(500, Error{
			Message:    err.Error(),
			StatusCode: 500,
			Code:       "500001",
		})
		return
	}
	fmt.Println(accountList)
	ctx.JSON(200, accountList)
}

func (api *AccountApi) CreateSession(ctx *gin.Context) {
	req, ok := ShouldBindJSON[CreateSessionRequest](ctx)
	if !ok {
		return
	}
	var account model.Account
	tx := api.db
	switch req.Idby {
	case "phone":
		var accountPhone model.AccountPhone
		err := api.db.First(&accountPhone, "namespace = ? AND phone_region = ? AND phone_number = ?", req.Namespace, req.PhoneRegion, req.PhoneNumber).Error
		if err != nil {
			ctx.JSON(404, Error{
				Message:    "account not found",
				StatusCode: 404,
				Code:       "404001",
			})
			return
		}
		tx = tx.Where("id = ?", accountPhone.ID)
	case "account":
		var accountTable model.AccountTable
		err := api.db.First(&accountTable, "namespace = ? AND account = ?", req.Namespace, req.Account).Error
		if err != nil {
			ctx.JSON(404, Error{
				Message:    "account not found",
				StatusCode: 404,
				Code:       "404001",
			})
			return
		}
		tx = tx.Where("id = ?", accountTable.ID)
	case "email":
		var accountEmail model.AccountEmail
		err := api.db.First(&accountEmail, "namespace = ? AND email = ?", req.Namespace, req.Email).Error
		if err != nil {
			ctx.JSON(404, Error{
				Message:    "account not found",
				StatusCode: 404,
				Code:       "404001",
			})
			return
		}
		tx = tx.Where("id = ?", accountEmail.ID)
	default:
		tx = tx.Where("id = ?", req.Account)
	}
	err := tx.First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, Error{
				Message:    "account not found",
				StatusCode: 404,
				Code:       "404001",
			})
			return
		}
		ctx.JSON(500, Error{
			Message:    err.Error(),
			StatusCode: 500,
			Code:       "500001",
		})
	}
	if !account.CheckPassword(req.Password) {
		ctx.JSON(401, Error{
			Message:    "password not match",
			StatusCode: 401,
			Code:       "401001",
		})
		return
	}
	token, _ := SignJWT(jwt.MapClaims{
		"type":      "admin",
		"namespace": account.Namespace,
		"account":   account.Account,
		"email":     account.Email,
		"id":        account.ID,
		"exp":       time.Now().Add(time.Hour * 7 * 24).Unix(),
	})
	ctx.JSON(200, gin.H{
		"token": token,
	})
}

func (api *AccountApi) UpdatePassword(c *gin.Context) {
	account, ok := GetAccount(api.db, c)
	if !ok {
		c.JSON(404, Error{
			Message:    "account not found",
			StatusCode: 404,
			Code:       "404001",
		})
		return
	}
	req, ok := ShouldBindJSON[SetPasswordRequest](c)
	if !ok {
		return
	}
	account.SetPassword(req.Password)
	err := api.db.Save(account).Error
	if err != nil {
		c.JSON(500, Error{})
	}
	c.JSON(204, nil)
}

func (api *AccountApi) CreateAccount(c *gin.Context) {
	req, ok := ShouldBindJSON[CreateAccountRequest](c)
	if !ok {
		return
	}
	account := model.Account{
		Namespace: req.Namespace,
	}
	if req.ID != "" {
		account.ID = req.ID
	} else {
		node, _ := snowflake.NewNode(0)
		account.ID = node.Generate().String()
	}
	if req.Account != "" {
		account.Account = &req.Account
	}
	if req.Email != "" {
		account.Email = &req.Email
	}
	if req.PhoneRegion != "" {
		account.PhoneRegion = &req.PhoneRegion
	}
	if req.PhoneNumber != "" {
		account.PhoneNumber = &req.PhoneNumber
	}
	if req.Password != "" {
		account.SetPassword(req.Password)
	}
	err := api.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&account).Error
		if err != nil {
			tx.Rollback()
		}
		return err
	})
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // MySQL code for duplicate entry
				c.JSON(409, Error{
					Message:    "key conflict. maybe accountId or account or email or phone number already exists",
					StatusCode: 409,
					Code:       "409001",
				})
			}
			return
		}
		c.JSON(500, Error{
			Message:    err.Error(),
			StatusCode: 500,
			Code:       "500001",
		})
		return
	}
	c.JSON(201, account)
}

func GetAccount(db *gorm.DB, ctx *gin.Context) (*model.Account, bool) {
	namespace := ctx.Query("namespace")
	idby := ctx.Query("idby")
	id := ctx.Query("id")
	accountString := ctx.Query("account")
	email := ctx.Query("email")
	var account model.Account
	tx := db
	switch idby {
	case "phone":
		var accountPhone model.AccountPhone
		err := db.First(&accountPhone, "namespace = ? AND phone_region = ? AND phone_number = ?", namespace, id, email).Error
		if err != nil {
			return nil, false
		}
		tx = tx.Where("id = ?", accountPhone.ID)
	case "account":
		var accountTable model.AccountTable
		err := db.First(&accountTable, "namespace = ? AND account = ?", namespace, accountString).Error
		if err != nil {
			return nil, false
		}
		tx = tx.Where("id = ?", accountTable.ID)
	case "email":
		var accountEmail model.AccountEmail
		err := db.First(&accountEmail, "namespace = ? AND email = ?", namespace, email).Error
		if err != nil {
			return nil, false
		}
		tx = tx.Where("id = ?", accountEmail.ID)
	default:
		tx = tx.Where("namespace = ? AND id = ?", namespace, id)
	}
	err := tx.First(&account).Error
	if err != nil {
		return nil, false
	}
	return &account, true
}
