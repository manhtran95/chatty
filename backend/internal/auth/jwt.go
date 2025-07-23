package auth

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, expiry time.Duration) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateAccessToken(userID string) string {
	exp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRE_MINUTES"))
	accessTokenExpiry := time.Minute * time.Duration(exp)
	token, err := GenerateToken(userID, accessTokenExpiry)
	if err != nil {
		log.Fatal("Error generating access token:", err)
	}
	return token
}

func GenerateRefreshToken(userID string) string {
	exp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRE_DAYS"))
	refreshTokenExpiry := time.Hour * 24 * time.Duration(exp)
	token, err := GenerateToken(userID, refreshTokenExpiry)
	if err != nil {
		log.Fatal("Error generating refresh token:", err)
	}
	return token
}

func VerifyAndParseToken(tokenStr string) (*jwt.Token, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
}

func ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := VerifyAndParseToken(tokenStr)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", jwt.ErrSignatureInvalid
}
