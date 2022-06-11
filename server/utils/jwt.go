package utils

import (
	"os"
	"time"

	// gojwt "github.com/dgrijalva/jwt-go"
	// gojwt "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

func init() {
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
}

type Claims struct {
	Uid int
	jwt.RegisteredClaims
}

func Award(uid int) (string, error) {
	// 过期时间 默认7天
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtKey)
	return token, err
}

func ParseToken(tokenStr string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}
