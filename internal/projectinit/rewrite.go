package projectinit

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const secretBytes = 32

var skipDirNames = map[string]struct{}{
	".git":         {},
	".idea":        {},
	".vscode":      {},
	".umi":         {},
	".turbopack":   {},
	"node_modules": {},
	"vendor":       {},
	"bin":          {},
	"dist":         {},
	"tmp":          {},
	"temp":         {},
	"storage":      {},
}

var skipFileNames = map[string]struct{}{
	"go.sum":            {},
	"pnpm-lock.yaml":    {},
	"pnpm-lock.yml":     {},
	"package-lock.json": {},
	"yarn.lock":         {},
}

var skipExts = map[string]struct{}{
	".png":   {},
	".jpg":   {},
	".jpeg":  {},
	".gif":   {},
	".webp":  {},
	".ico":   {},
	".woff":  {},
	".woff2": {},
	".ttf":   {},
	".eot":   {},
	".so":    {},
	".exe":   {},
	".dll":   {},
	".zip":   {},
	".gz":    {},
	".wasm":  {},
	".db":    {},
	".out":   {},
}

// Result 是一次清名扫描的结果。
type Result struct {
	Files     []string
	JWTSecret string
	CryptoKey string
	DryRun    bool
}

// Secrets 保存 init 生成的随机密钥，测试可注入以断言写回值。
type Secrets struct {
	JWT    string
	Crypto string
}

// Run 在 root 下按身份替换模板残留名。
// dryRun 只收集将改路径，不写盘。
func Run(root string, id Identity, secrets Secrets, dryRun bool) (*Result, error) {
	id, err := id.Normalize()
	if err != nil {
		return nil, err
	}
	if err := secrets.ensure(); err != nil {
		return nil, err
	}

	root, err = filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("解析项目根目录失败: %w", err)
	}

	result := &Result{
		JWTSecret: secrets.JWT,
		CryptoKey: secrets.Crypto,
		DryRun:    dryRun,
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		name := d.Name()
		if d.IsDir() {
			if _, skip := skipDirNames[name]; skip && path != root {
				return fs.SkipDir
			}
			return nil
		}
		if shouldSkipFile(name) {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		changed, err := rewriteFile(path, rel, id, secrets, dryRun)
		if err != nil {
			return err
		}
		if changed {
			result.Files = append(result.Files, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GenerateSecrets 生成 JWT 与前端记住密码使用的独立随机密钥。
func GenerateSecrets() (Secrets, error) {
	jwt, err := randomHex(secretBytes)
	if err != nil {
		return Secrets{}, err
	}
	crypto, err := randomHex(secretBytes)
	if err != nil {
		return Secrets{}, err
	}
	return Secrets{JWT: jwt, Crypto: crypto}, nil
}

// ensure 在调用方未注入密钥时生成，保证配置写回的是随机值而不是项目名拼接。
func (s *Secrets) ensure() error {
	if s.JWT != "" && s.Crypto != "" {
		return nil
	}
	generated, err := GenerateSecrets()
	if err != nil {
		return err
	}
	if s.JWT == "" {
		s.JWT = generated.JWT
	}
	if s.Crypto == "" {
		s.Crypto = generated.Crypto
	}
	return nil
}

// shouldSkipFile 跳过锁文件、校验和与静态/二进制产物，避免无意义的大文件改写。
func shouldSkipFile(name string) bool {
	if _, skip := skipFileNames[name]; skip {
		return true
	}
	ext := strings.ToLower(filepath.Ext(name))
	_, skip := skipExts[ext]
	return skip
}

// rewriteFile 读取单个文本文件并按规则改写。
// 含 NUL 的文件视为产物，直接跳过。
func rewriteFile(absPath, rel string, id Identity, secrets Secrets, dryRun bool) (bool, error) {
	info, err := os.Stat(absPath)
	if err != nil {
		return false, err
	}
	if info.Size() > 8<<20 {
		return false, nil
	}

	raw, err := os.ReadFile(absPath)
	if err != nil {
		return false, err
	}
	if bytes.IndexByte(raw, 0) >= 0 {
		return false, nil
	}

	next := apply(string(raw), rel, id, secrets)
	if next == string(raw) {
		return false, nil
	}
	if dryRun {
		return true, nil
	}
	if err := os.WriteFile(absPath, []byte(next), info.Mode().Perm()); err != nil {
		return false, fmt.Errorf("写入 %s 失败: %w", rel, err)
	}
	return true, nil
}

// apply 按优先级改写文件内容。
// JWT / crypto 密钥走专用规则，避免把默认 secret 当成普通项目名替换。
func apply(content, rel string, id Identity, secrets Secrets) string {
	if isAppConfig(rel) {
		content = strings.ReplaceAll(content, oldJWTSecret, secrets.JWT)
	}
	if filepath.Base(rel) == "crypto.ts" {
		content = strings.ReplaceAll(content, oldCryptoSecret, secrets.Crypto)
	}
	for _, rule := range id.replacements() {
		content = strings.ReplaceAll(content, rule.old, rule.new)
	}
	return content
}

// isAppConfig 只在应用配置里写入新 JWT，避免文档示例被替换成真实密钥。
func isAppConfig(rel string) bool {
	ext := strings.ToLower(filepath.Ext(rel))
	if ext != ".yaml" && ext != ".yml" {
		return false
	}
	dir := filepath.ToSlash(filepath.Dir(rel))
	base := filepath.Base(rel)
	if base == "config.yaml" || base == "config.yml" {
		return true
	}
	return dir == "config" || strings.HasSuffix(dir, "/config")
}

// randomHex 生成指定字节数的十六进制密钥，满足 JWT 至少 32 位的校验。
func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("生成随机密钥失败: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
