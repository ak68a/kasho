package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTToken struct {
	config *Config
}

type jwtClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
	Exp int64 `json:"exp"`
}

func NewJWTToken(config *Config) *JWTToken {
	return &JWTToken{config: config}
}

func (j * JWTToken) CreateToken(userID int64) (string, error) {
	claims := jwtClaim{
		UserID: userID,
		// Exp: time.Now().Add(time.Minute * 30).Unix(),
		// Exp: time.Now().Add(time.Hour * 24).Unix(),
		Exp: time.Now().Add(time.Second * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.config.Signing_key))
	if err != nil {
		return "", err
	}

	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error) {
	fmt.Printf("Received token: %s\n", tokenString)
	
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Invalid signing method: %v\n", t.Method)
			return nil, fmt.Errorf("invalid token")
		}
		fmt.Printf("Using signing key: %s\n", j.config.Signing_key)
		return []byte(j.config.Signing_key), nil
	})

	if err != nil {
		fmt.Printf("Token parsing error: %v\n", err)
		return 0, err
	}

	claims, ok := token.Claims.(*jwtClaim)
	if !ok {
		fmt.Println("Failed to parse claims")
		return 0, fmt.Errorf("invalid token")
	}
	if !token.Valid {
		fmt.Println("Token is invalid")
		return 0, fmt.Errorf("invalid token")
	}

	if claims.Exp < time.Now().Unix() {
		fmt.Printf("Token expired at %v, current time: %v\n", claims.Exp, time.Now().Unix())
		return 0, fmt.Errorf("token expired")
	}

	fmt.Printf("Token verified successfully for user ID: %v\n", claims.UserID)
	return claims.UserID, nil
}