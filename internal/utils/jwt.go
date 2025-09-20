package utils

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/LamaKhaledd/HeartReach/internal/db"
)

func CreateJwtAccessToken(key string, user db.User) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "heartreach-access",
		Subject:   fmt.Sprintf("%v", user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func GetDataFromToken(token, key string) (string, string, error) {
	tokenInterface, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", "", err
	}

	if claims, ok := tokenInterface.Claims.(*jwt.RegisteredClaims); ok && tokenInterface.Valid {
		userId := claims.Subject
		issuer := claims.Issuer
		return userId, issuer, nil
	} else {
		return "", "", err
	}
}
