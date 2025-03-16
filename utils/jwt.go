package utils

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

var JwtSecret = []byte("by_kq")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 令牌有效期为24小时
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	//封装
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//密钥（依据这个加密）
	return token.SignedString(JwtSecret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func GetName(c context.Context, ctx *app.RequestContext) {
	// 从上下文中获取用户名
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, utils.H{"error": "Username not found"})
		return
	}

	ctx.JSON(http.StatusOK, utils.H{
		"message":  "Welcome!",
		"username": username,
	})
}
