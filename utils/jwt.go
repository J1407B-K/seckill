package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var JwtSecret = []byte("by_kq")

type Claims struct {
	UserId string `json:"UserId"`
	jwt.RegisteredClaims
}

func GenerateToken(userid string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 令牌有效期为24小时
	claims := &Claims{
		UserId: userid,
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
