package projectinit

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	// 旧 Go module，用于改写 go.mod 与 import 路径。
	oldModulePath = "bico-admin"
	// 旧 GitHub 仓库地址，文档与 swagger 按新 module 推导。
	oldGitHubPath = "github.com/slowlyo/bico-admin"
	// 旧前端 package.json name。
	oldFrontendName = "ant-design-pro"
	// 旧前端仓库地址，需在通用 ant-design-pro 替换之前处理，避免改成错误 host。
	oldFrontendRepo = "git@github.com:ant-design/ant-design-pro.git"
	// 旧展示名。
	oldTitle = "Bico Admin"
	// 旧数据库名。
	oldDatabase = "bico_admin"
	// 旧环境变量前缀（不含尾部下划线）。
	oldEnvPrefix = "BICO"
	// 配置里需要被随机密钥替换的默认 JWT 值，避免只做项目名拼接。
	oldJWTSecret = "bico-admin-secret-key-change-in-production"
	// 前端记住密码用的硬编码密钥。
	oldCryptoSecret = "bico-admin-secret-key"
)

// Identity 描述清名后的新项目身份。
type Identity struct {
	Module string
	Name   string
	Title  string
}

// replacement 是按优先级排列的一对替换规则。
type replacement struct {
	old string
	new string
}

// Normalize 补齐可推导字段并校验。
// name 缺省时取 module 最后一段；title 缺省时由 name 转成空格分隔的标题。
func (id Identity) Normalize() (Identity, error) {
	id.Module = strings.TrimSpace(id.Module)
	id.Name = strings.TrimSpace(id.Name)
	id.Title = strings.TrimSpace(id.Title)

	if id.Module == "" {
		return Identity{}, fmt.Errorf("必须指定 --module")
	}
	if strings.ContainsAny(id.Module, " \t\n") || strings.HasPrefix(id.Module, "/") || strings.HasSuffix(id.Module, "/") {
		return Identity{}, fmt.Errorf("module 格式无效: %s", id.Module)
	}

	if id.Name == "" {
		id.Name = id.Module
		if i := strings.LastIndex(id.Name, "/"); i >= 0 {
			id.Name = id.Name[i+1:]
		}
	}
	if err := validateName(id.Name); err != nil {
		return Identity{}, err
	}

	if id.Title == "" {
		id.Title = titleFromName(id.Name)
	}
	return id, nil
}

// SnakeName 把二进制名转成数据库风格的下划线名。
func (id Identity) SnakeName() string {
	return strings.ReplaceAll(id.Name, "-", "_")
}

// EnvPrefix 生成 Viper / Compose 使用的大写前缀，对应原 BICO。
func (id Identity) EnvPrefix() string {
	return strings.ToUpper(id.SnakeName())
}

// GitSSH 按新 module 推导 package.json 的 git 仓库地址。
func (id Identity) GitSSH() string {
	host, rest, ok := strings.Cut(id.Module, "/")
	if !ok {
		return id.Module
	}
	return "git@" + host + ":" + rest + ".git"
}

// replacements 返回替换表，顺序不可改：先处理完整路径，再处理短名，避免 import 被写成二进制名。
func (id Identity) replacements() []replacement {
	return []replacement{
		{oldGitHubPath, id.Module},
		{`"` + oldModulePath + `/`, `"` + id.Module + `/`},
		{"`" + oldModulePath + "/", "`" + id.Module + "/"},
		{"module " + oldModulePath, "module " + id.Module},
		{oldFrontendRepo, id.GitSSH()},
		{oldTitle, id.Title},
		{oldFrontendName, id.Name},
		{oldModulePath, id.Name},
		{oldDatabase, id.SnakeName()},
		{"bico-mysql", id.Name + "-mysql"},
		{"bico-redis", id.Name + "-redis"},
		{"bico-network", id.Name + "-network"},
		{oldEnvPrefix, id.EnvPrefix()},
	}
}

// validateName 限制二进制 / 镜像 / 容器名使用小写短横线，避免 Docker 与 Makefile 目标非法。
func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("必须指定 --name")
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			continue
		case r >= '0' && r <= '9':
			if i == 0 {
				return fmt.Errorf("name 必须以字母开头: %s", name)
			}
		case r == '-':
			if i == 0 || i == len(name)-1 {
				return fmt.Errorf("name 不能以短横线开头或结尾: %s", name)
			}
		default:
			return fmt.Errorf("name 仅允许小写字母、数字与短横线: %s", name)
		}
	}
	return nil
}

// titleFromName 把 foo-admin 转成 Foo Admin，供缺省 --title 使用。
func titleFromName(name string) string {
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, " ")
}
