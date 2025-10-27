package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AnimalClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func getSecretKey() ([]byte, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("FATAL: JWT_SECRET not found. Please check your .env file.")
    }
    return []byte(secret), nil
}

func GenerateTokenPair(userID uint) (accessToken string, refreshToken string, err error) {
    jwtSecret, err := getSecretKey()
    if err != nil {
        return "", "", err
    }

    accessExpMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"))
    refreshExpHours, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_HOURS"))

    accessExpiry := time.Now().Add(time.Duration(accessExpMinutes) * time.Minute)
    accessClaims := &AnimalClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExpiry),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Audience:  jwt.ClaimStrings{"access"}, 
        },
    }
    accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtSecret)
    if err != nil {
        return "", "", err
    }

    refreshExpiry := time.Now().Add(time.Duration(refreshExpHours) * time.Hour)
    refreshClaims := &AnimalClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(refreshExpiry),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Audience:  jwt.ClaimStrings{"refresh"}, 
        },
    }
    refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtSecret)
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
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