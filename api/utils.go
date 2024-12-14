package api

import (
	"github.com/gin-gonic/gin"
	"github.com/macaoservices/auth/fault"
)

func ShouldBindJSON[T any](ctx *gin.Context) (T, bool) {
	var t T
	err := ctx.ShouldBindJSON(&t)
	if err != nil {
		fault.ErrBadRequest.GinHandler(ctx)
		return t, false
	}
	return t, true
}
