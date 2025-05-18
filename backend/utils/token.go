package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)


type jwtClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
	Exp int64 `json:"exp"`
}

func CreateToken(userID int64, signingKey string) (string, error) {
	claims := jwtClaim{
		UserID: userID,
		Exp: time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return string(tokenString), nil
}

func VerifyToken(tokenString, signingKey string) (int64, error) {
	token, err :=jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwtClaim)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("Invalid token")
	}

	if claims.Exp < time.Now().Unix() {
		return 0, fmt.Errorf("Token expired")
	}

	return claims.UserID, nil
}