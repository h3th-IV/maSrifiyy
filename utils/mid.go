package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maSrifiyy/models"
)

var (
	SECRET = "thugnificient@booming.metro"
	ISSUER = "g-unit"
)

func GenerateJWT(seller models.Sellers, exp time.Duration, issuer, secret string) (string, error) {
	var expData = time.Now().Add(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   seller.UserID,
		"exp":    expData.Unix(),
		"issuer": issuer,
	})

	JWToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return JWToken, nil
}

func DecodeJWT(tokenString, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to parse token claims")
	}

	return claims, nil
}
