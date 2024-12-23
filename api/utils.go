package api

import "github.com/gin-gonic/gin"

func ShouldBindJSON[T any](ctx *gin.Context) (T, bool) {
	var t T
	err := ctx.ShouldBindJSON(&t)
	if err != nil {
		ctx.JSON(400, Error{
			Message:    "bad request",
			Code:       "400001",
			StatusCode: 400,
		})
		return t, false
	}
	return t, true
}
