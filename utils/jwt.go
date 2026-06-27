package utils

import (
	"blog/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string) (string, error) {
	// 从环境变量获取jwt secret
	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key")
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // 签发日期
			NotBefore: jwt.NewNumericDate(time.Now()),                     // 有效起始日期
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// 从环境变量获取jwt secret
	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
