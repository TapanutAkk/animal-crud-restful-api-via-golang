package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AnimalClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func getSecretKey() ([]byte, error) {
    // ดึงค่า SECRET key จาก Environment (ซึ่งถูกโหลดแล้วโดย main.go init())
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        // ใช้ log.Fatal จะหยุดโปรแกรมทันที
        log.Fatal("FATAL: JWT_SECRET not found. Please check your .env file.")
        // หรือ return error 
        // return nil, fmt.Errorf("JWT_SECRET is missing")
    }
    return []byte(secret), nil
}

func GenerateToken(userID uint) (string, error) {
	jwtSecret, err := getSecretKey()
    if err != nil {
        return "", err
    }

	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &AnimalClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

func ValidateToken(tokenString string) (*AnimalClaims, error) {
	jwtSecret, err := getSecretKey()
    if err != nil {
        return nil, err
    }
	
	claims := &AnimalClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}