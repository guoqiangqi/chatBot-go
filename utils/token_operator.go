package chatbot

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID string, secretKey []byte) (string, error) {
	expireTime := time.Now().Add(time.Hour * 24).Unix()

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "chatbot",
			Subject:   useridTable[userID].Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
