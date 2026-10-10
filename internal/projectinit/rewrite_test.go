package projectinit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRewriteReplacesModuleImportsAndBinaryName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module bico-admin\n\ngo 1.25\n")
	writeFile(t, root, "cmd/main.go", "package main\nimport \"bico-admin/internal/api\"\nvar root = \"bico-admin\"\n")
	writeFile(t, root, "Makefile", "\t@go build -o bin/bico-admin ./cmd/main.go\n")

	result := runInit(t, root, false)

	goMod := readFile(t, root, "go.mod")
	if !strings.Contains(goMod, "module github.com/acme/foo") {
		t.Fatalf("go.mod 未改写 module: %s", goMod)
	}
	mainGo := readFile(t, root, "cmd/main.go")
	if !strings.Contains(mainGo, `"github.com/acme/foo/internal/api"`) {
		t.Fatalf("import 未改成新 module: %s", mainGo)
	}
	if !strings.Contains(mainGo, `var root = "foo"`) {
		t.Fatalf("二进制名未改写: %s", mainGo)
	}
	makefile := readFile(t, root, "Makefile")
	if !strings.Contains(makefile, "bin/foo") || strings.Contains(makefile, "bico-admin") {
		t.Fatalf("Makefile 二进制名未改写: %s", makefile)
	}
	assertChanged(t, result, "go.mod", "cmd/main.go", "Makefile")
}

func TestRewriteJWTSecretIsRandomNotRenamed(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "config/config.yaml", "app:\n  name: Bico Admin\njwt:\n  secret: \"bico-admin-secret-key-change-in-production\"\n")
	writeFile(t, root, "internal/core/config/config.go", "if secret == \"bico-admin-secret-key-change-in-production\" {\n}\nv.SetEnvPrefix(\"BICO\")\n")
	writeFile(t, root, "docs/config.md", "  secret: \"bico-admin-secret-key-change-in-production\"\n")
	writeFile(t, root, "web/src/utils/crypto.ts", "const SECRET_KEY = 'bico-admin-secret-key';\n")

	result := runInit(t, root, false)

	cfg := readFile(t, root, "config/config.yaml")
	if strings.Contains(cfg, oldJWTSecret) || strings.Contains(cfg, "foo-secret-key-change-in-production") {
		t.Fatalf("jwt.secret 仍是字符串替换结果: %s", cfg)
	}
	if !strings.Contains(cfg, result.JWTSecret) {
		t.Fatalf("jwt.secret 未写入随机值")
	}
	if len(result.JWTSecret) < 32 {
		t.Fatalf("jwt.secret 长度不足: %d", len(result.JWTSecret))
	}
	if strings.Contains(cfg, "Bico Admin") {
		t.Fatalf("app.name 未替换")
	}

	code := readFile(t, root, "internal/core/config/config.go")
	if strings.Contains(code, result.JWTSecret) {
		t.Fatalf("校验代码不应写入真实密钥")
	}
	if !strings.Contains(code, `SetEnvPrefix("FOO")`) {
		t.Fatalf("环境变量前缀未替换: %s", code)
	}

	docs := readFile(t, root, "docs/config.md")
	if strings.Contains(docs, result.JWTSecret) {
		t.Fatalf("文档示例不应写入真实密钥: %s", docs)
	}

	crypto := readFile(t, root, "web/src/utils/crypto.ts")
	if strings.Contains(crypto, oldCryptoSecret) || strings.Contains(crypto, oldJWTSecret) {
		t.Fatalf("crypto.ts 仍是旧密钥")
	}
	if !strings.Contains(crypto, result.CryptoKey) {
		t.Fatalf("crypto.ts 未写入新密钥")
	}
}

func TestDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	original := "module bico-admin\n"
	writeFile(t, root, "go.mod", original)

	result := runInit(t, root, true)
	if len(result.Files) == 0 {
		t.Fatal("dry-run 应列出将改路径")
	}
	if readFile(t, root, "go.mod") != original {
		t.Fatal("dry-run 不应写盘")
	}
	if !result.DryRun {
		t.Fatal("结果应标记 dry-run")
	}
}

func TestSkipLockfilesAndChecksums(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module bico-admin\n")
	writeFile(t, root, "go.sum", "bico-admin h1:abc\n")
	writeFile(t, root, "web/pnpm-lock.yaml", "name: ant-design-pro\n")
	writeFile(t, root, "bin/bico-admin", "bico-admin")

	runInit(t, root, false)

	if readFile(t, root, "go.sum") != "bico-admin h1:abc\n" {
		t.Fatal("不应改写 go.sum")
	}
	if readFile(t, root, "web/pnpm-lock.yaml") != "name: ant-design-pro\n" {
		t.Fatal("不应改写 pnpm-lock")
	}
}

func TestRewriteDockerAndFrontendIdentity(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docker-compose.yml", "image: bico-admin:latest\ncontainer_name: bico-admin\nBICO_JWT_SECRET: x\nMYSQL_DATABASE: bico_admin\nnetworks:\n  - bico-network\n  bico-network:\ncontainer_name: bico-mysql\n")
	writeFile(t, root, "web/package.json", "{\n  \"name\": \"ant-design-pro\",\n  \"repository\": \"git@github.com:ant-design/ant-design-pro.git\"\n}\n")
	writeFile(t, root, "web/config/defaultSettings.ts", "title: 'Bico Admin',\n")
	writeFile(t, root, "docs/api/doc.go", "// @title Bico Admin Open API\n// @termsOfService https://github.com/slowlyo/bico-admin\n")

	runInit(t, root, false)

	compose := readFile(t, root, "docker-compose.yml")
	for _, leftover := range []string{"bico-admin", "bico_admin", "BICO_", "bico-network", "bico-mysql"} {
		if strings.Contains(compose, leftover) {
			t.Fatalf("compose 仍有 %s: %s", leftover, compose)
		}
	}
	pkg := readFile(t, root, "web/package.json")
	if strings.Contains(pkg, "ant-design-pro") {
		t.Fatalf("package.json 仍有 ant-design-pro: %s", pkg)
	}
	if !strings.Contains(pkg, "git@github.com:acme/foo.git") {
		t.Fatalf("package.json 仓库地址未按 module 推导: %s", pkg)
	}
	settings := readFile(t, root, "web/config/defaultSettings.ts")
	if !strings.Contains(settings, "Foo Admin") {
		t.Fatalf("title 未替换: %s", settings)
	}
	doc := readFile(t, root, "docs/api/doc.go")
	if !strings.Contains(doc, "Foo Admin Open API") || !strings.Contains(doc, "https://github.com/acme/foo") {
		t.Fatalf("swagger 注释未替换: %s", doc)
	}
}

func TestNormalizeDerivesNameAndTitle(t *testing.T) {
	id, err := Identity{Module: "github.com/acme/foo-admin"}.Normalize()
	if err != nil {
		t.Fatalf("Normalize 失败: %v", err)
	}
	if id.Name != "foo-admin" {
		t.Fatalf("未从 module 推导 name: %s", id.Name)
	}
	if id.Title != "Foo Admin" {
		t.Fatalf("未从 name 推导 title: %s", id.Title)
	}
	if id.EnvPrefix() != "FOO_ADMIN" {
		t.Fatalf("环境变量前缀错误: %s", id.EnvPrefix())
	}
}

func TestNormalizeRejectsInvalidName(t *testing.T) {
	_, err := Identity{Module: "github.com/acme/foo", Name: "Foo"}.Normalize()
	if err == nil {
		t.Fatal("大写 name 应被拒绝")
	}
}

func runInit(t *testing.T, root string, dryRun bool) *Result {
	t.Helper()
	result, err := Run(root, Identity{
		Module: "github.com/acme/foo",
		Name:   "foo",
		Title:  "Foo Admin",
	}, Secrets{}, dryRun)
	if err != nil {
		t.Fatalf("Run 失败: %v", err)
	}
	return result
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("创建目录失败: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("写入 %s 失败: %v", rel, err)
	}
}

func readFile(t *testing.T, root, rel string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("读取 %s 失败: %v", rel, err)
	}
	return string(raw)
}

func assertChanged(t *testing.T, result *Result, files ...string) {
	t.Helper()
	found := map[string]struct{}{}
	for _, file := range result.Files {
		found[file] = struct{}{}
	}
	for _, file := range files {
		if _, ok := found[file]; !ok {
			t.Fatalf("缺少变更文件 %s, 实际: %v", file, result.Files)
		}
	}
}
