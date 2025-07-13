package jwt

import (
	"bico-admin/pkg/cache"
	"context"
	"errors"
	"fmt"
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
	cache      cache.Cache
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret, issuer string, expireTime time.Duration) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		issuer:     issuer,
		expireTime: expireTime,
	}
}

// NewJWTManagerWithCache 创建带缓存的JWT管理器
func NewJWTManagerWithCache(secret, issuer string, expireTime time.Duration, cache cache.Cache) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		issuer:     issuer,
		expireTime: expireTime,
		cache:      cache,
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

// AddToBlacklist 将令牌加入黑名单
func (j *JWTManager) AddToBlacklist(ctx context.Context, tokenString string) error {
	if j.cache == nil {
		return errors.New("缓存客户端未配置")
	}

	// 解析令牌获取过期时间
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return fmt.Errorf("解析令牌失败: %w", err)
	}

	// 计算令牌剩余有效时间
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration <= 0 {
		// 令牌已过期，无需加入黑名单
		return nil
	}

	// 将令牌加入缓存黑名单，设置过期时间为令牌的剩余有效时间
	key := fmt.Sprintf("jwt_blacklist:%s", tokenString)
	return j.cache.Set(ctx, key, "1", expiration)
}

// IsBlacklisted 检查令牌是否在黑名单中
func (j *JWTManager) IsBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	if j.cache == nil {
		// 如果没有缓存，则认为令牌不在黑名单中
		return false, nil
	}

	key := fmt.Sprintf("jwt_blacklist:%s", tokenString)
	exists, err := j.cache.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("检查黑名单失败: %w", err)
	}

	return exists, nil
}
