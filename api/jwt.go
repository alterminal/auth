package api

import (
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(viper.GetString("jwt.secret"))
}

func SignJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString, nil
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	return token.Claims.(jwt.MapClaims), err
}
