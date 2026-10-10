package jwt

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestTokenRoundTrip 验证生成的 token 可被同一密钥解析，且载荷完整。
func TestTokenRoundTrip(t *testing.T) {
	m := NewJWTManager("secret", 1)
	token, err := m.GenerateToken(7, "alice", 3)
	if err != nil {
		t.Fatalf("生成 token 失败: %v", err)
	}
	claims, err := m.ParseToken(token)
	if err != nil {
		t.Fatalf("解析 token 失败: %v", err)
	}
	if claims.UserID != 7 || claims.Username != "alice" || claims.Version != 3 {
		t.Fatalf("载荷不一致: %+v", claims)
	}
}

// TestParseTokenRejectsInvalid 验证签名被篡改、密钥不符、结构错误的 token 都返回无效。
func TestParseTokenRejectsInvalid(t *testing.T) {
	m := NewJWTManager("secret", 1)
	token, _ := m.GenerateToken(1, "alice", 1)
	parts := strings.Split(token, ".")

	// 把载荷换成另一个用户但保留原签名，签名校验必须失败。
	other, _ := m.GenerateToken(2, "mallory", 1)
	forged := parts[0] + "." + strings.Split(other, ".")[1] + "." + parts[2]

	cases := map[string]string{
		"签名被篡改": parts[0] + "." + parts[1] + "." + parts[2] + "x",
		"载荷被替换": forged,
		"缺少签名段": parts[0] + "." + parts[1],
		"空字符串":  "",
		"段数过多":  token + ".extra",
	}
	for name, tk := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := m.ParseToken(tk); !errors.Is(err, ErrTokenInvalid) {
				t.Fatalf("期望 ErrTokenInvalid，实际 %v", err)
			}
		})
	}

	// 不同密钥签发的 token 不能被当前密钥接受。
	if _, err := NewJWTManager("another", 1).ParseToken(token); !errors.Is(err, ErrTokenInvalid) {
		t.Fatalf("密钥不符应无效，实际 %v", err)
	}
}

// TestParseTokenExpired 验证签名正确但已过期的 token 返回过期错误。
func TestParseTokenExpired(t *testing.T) {
	expired := createToken(&Claims{UserID: 1, Exp: time.Now().Add(-time.Minute).Unix()}, "secret")
	if _, err := NewJWTManager("secret", 1).ParseToken(expired); !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("期望 ErrTokenExpired，实际 %v", err)
	}
}
