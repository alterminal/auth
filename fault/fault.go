package fault

import "github.com/gin-gonic/gin"

type Error struct {
	StatusCode int
	ErrorCode  string
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GinHandler(c *gin.Context) {
	c.JSON(e.StatusCode, gin.H{"message": e.Message, "code": e.ErrorCode})
}

var ErrNotFound = &Error{
	StatusCode: 404,
	ErrorCode:  "404001",
	Message:    "Not Found",
}

var ErrBadRequest = &Error{
	StatusCode: 400,
	ErrorCode:  "400001",
	Message:    "Bad Request",
}

var ErrUnauthorized = &Error{
	StatusCode: 401,
	ErrorCode:  "401001",
	Message:    "Unauthorized",
}


