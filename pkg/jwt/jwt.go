package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	secret     []byte
	issuer     string
	expireTime time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret, issuer string, expireTime time.Duration) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		issuer:     issuer,
		expireTime: expireTime,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWTManager) GenerateToken(userID uint, username, userType string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.expireTime)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken 验证JWT令牌
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新令牌
func (j *JWTManager) RefreshToken(tokenString string) (string, time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	// 检查令牌是否即将过期（剩余时间少于1小时）
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", time.Time{}, errors.New("令牌尚未到刷新时间")
	}

	// 生成新令牌
	return j.GenerateToken(claims.UserID, claims.Username, claims.UserType)
}
