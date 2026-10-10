package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func createToken(claims *Claims, secret string) string {
	h := header{Alg: "HS256", Typ: "JWT"}

	headerJSON, _ := json.Marshal(h)
	claimsJSON, _ := json.Marshal(claims)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	message := headerB64 + "." + claimsB64
	signature := sign(message, secret)

	return message + "." + signature
}

func parseToken(tokenString, secret string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrTokenInvalid
	}

	message := parts[0] + "." + parts[1]
	signature := parts[2]

	// 使用常量时间比较签名，避免通过响应耗时逐字节猜出合法签名。
	if !hmac.Equal([]byte(sign(message, secret)), []byte(signature)) {
		return nil, ErrTokenInvalid
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrTokenInvalid
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, ErrTokenInvalid
	}

	return &claims, nil
}

func sign(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
