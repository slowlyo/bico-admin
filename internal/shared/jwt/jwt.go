package jwt

import (
	"errors"
	"time"
)

var (
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenInvalid = errors.New("token 无效")
)

// Claims JWT 载荷
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// GetExpiration 获取过期时间
func (c *Claims) GetExpiration() int64 {
	return c.Exp
}

// JWTManager JWT 管理器
type JWTManager struct {
	secret      string
	expireHours int
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secret string, expireHours int) *JWTManager {
	return &JWTManager{
		secret:      secret,
		expireHours: expireHours,
	}
}

// GenerateToken 生成 token
func (j *JWTManager) GenerateToken(userID uint, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Exp:      time.Now().Add(time.Duration(j.expireHours) * time.Hour).Unix(),
	}
	
	token := createToken(claims, j.secret)
	return token, nil
}

// ParseToken 解析 token
func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString, j.secret)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	
	if time.Now().Unix() > claims.Exp {
		return nil, ErrTokenExpired
	}
	
	return claims, nil
}

// ValidateToken 验证 token 并返回 map 格式的 claims
func (j *JWTManager) ValidateToken(tokenString string) (map[string]interface{}, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"user_id":  float64(claims.UserID),
		"username": claims.Username,
		"exp":      float64(claims.Exp),
	}, nil
}
