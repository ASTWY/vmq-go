package jwt

import (
	"time"
	"vmq-go/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// SecretKey 密钥
	secretKey = []byte(config.Conf.Jwt.Secret)
	// ExpireTime 过期时间 秒
	expireTime = config.Conf.Jwt.Expire
)

// GenerateToken 生成JWT
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(expireTime) * time.Second)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,              // 用户名
		"exp":      expirationTime.Unix(), // 过期时间
	})
	// 生成token
	token, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// IsTokenValid 校验Token是否有效
func IsTokenValid(tokenString string) bool {
	token, err := ParseToken(tokenString)
	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		return time.Now().Before(expirationTime)
	}

	return false
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string) (string, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 生成新的token
		claims["exp"] = time.Now().Add(time.Duration(expireTime) * time.Second).Unix()
		return GenerateToken(claims["username"].(string))
	}

	return "", nil
}

// 强制过期
func ForceExpire(tokenString string) bool {
	token, err := ParseToken(tokenString)
	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims["exp"] = time.Now().Add(-time.Second).Unix()
		return true
	}

	return false
}
