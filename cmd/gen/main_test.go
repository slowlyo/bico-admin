package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseFields 验证字段解析的缺省值与各类非法输入。
func TestParseFields(t *testing.T) {
	// 类型缺省为 string，标签缺省为列名。
	fields, err := parseFields("title,sort:int:排序,pinned:bool")
	if err != nil {
		t.Fatalf("合法字段应解析成功: %v", err)
	}
	if len(fields) != 3 {
		t.Fatalf("期望 3 个字段，实际 %d", len(fields))
	}
	if fields[0].Type != "string" || fields[0].Label != "title" || fields[0].GoName != "Title" {
		t.Fatalf("缺省值不正确: %+v", fields[0])
	}
	if fields[1].Type != "int" || fields[1].Label != "排序" {
		t.Fatalf("int 字段不正确: %+v", fields[1])
	}

	invalid := map[string]string{
		"空参数":    "",
		"保留列":    "enabled:bool",
		"保留列 id": "id:int",
		"未知类型":   "title:json",
		"重复字段":   "title,title",
		"非法列名":   "Title",
	}
	for name, arg := range invalid {
		t.Run(name, func(t *testing.T) {
			if _, err := parseFields(arg); err == nil {
				t.Fatalf("%q 应返回错误", arg)
			}
		})
	}
}

// TestNewSpec 验证命名推导和参数校验。
func TestNewSpec(t *testing.T) {
	spec, err := newSpec("article_tag", "文章标签", "title:string:标题")
	if err != nil {
		t.Fatalf("合法参数应成功: %v", err)
	}
	if spec.Pascal != "ArticleTag" || spec.Camel != "articleTag" || spec.Kebab != "article-tag" ||
		spec.Group != "/article-tags" || spec.Table != "article_tags" {
		t.Fatalf("命名推导不正确: %+v", spec)
	}

	for _, name := range []string{"", "Article", "article-tag", "_article", "article__tag"} {
		if _, err := newSpec(name, "文章", "title"); err == nil {
			t.Fatalf("模块名 %q 应被拒绝", name)
		}
	}
	if _, err := newSpec("article", "", "title"); err == nil {
		t.Fatal("缺少中文名称应被拒绝")
	}
}

// TestSpecColumns 验证搜索列只含 string，排序列排除 text。
func TestSpecColumns(t *testing.T) {
	spec, _ := newSpec("article", "文章", "title:string,body:text,sort:int,pinned:bool")
	if got := len(spec.SearchFields()); got != 1 || spec.SearchFields()[0].Name != "title" {
		t.Fatalf("搜索字段应只有 title，实际 %+v", spec.SearchFields())
	}
	got := strings.Join(spec.SortableColumns(), ",")
	want := "id,title,sort,pinned,enabled,created_at,updated_at"
	if got != want {
		t.Fatalf("排序列期望 %s，实际 %s", want, got)
	}
}

// TestInsertBeforeAnchor 验证插入位置、缩进，以及缺少锚点时报错。
func TestInsertBeforeAnchor(t *testing.T) {
	file := filepath.Join(t.TempDir(), "x.go")
	if err := os.WriteFile(file, []byte("a\n\t// gen:x\nb\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := insertBeforeAnchor(file, "// gen:x", "one,\ntwo,"); err != nil {
		t.Fatalf("插入失败: %v", err)
	}
	data, _ := os.ReadFile(file)
	if string(data) != "a\n\tone,\n\ttwo,\n\t// gen:x\nb\n" {
		t.Fatalf("插入结果不正确: %q", data)
	}

	if err := insertBeforeAnchor(file, "// gen:missing", "x"); err == nil {
		t.Fatal("缺少锚点应返回错误")
	}
}

// setupFakeRepo 在临时目录构造带锚点的注册文件，并切换为当前目录。
func setupFakeRepo(t *testing.T) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"internal/admin/module.go":                           "return []crud.Module{\n\t\t// gen:modules\n\t}\n",
		"internal/migrate/migrate.go":                        "db.AutoMigrate(\n\t\t// gen:models\n\t)\n",
		"web/apps/admin/src/router/routes/modules/system.ts": "children: [\n      // gen:routes\n    ],\n",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(root)
}

// TestGenerate 验证生成全部文件、注册到锚点，且重复生成会中止并保持原样。
func TestGenerate(t *testing.T) {
	setupFakeRepo(t)
	spec, err := newSpec("article", "文章", "title:string:标题,sort:int:排序")
	if err != nil {
		t.Fatal(err)
	}
	if err := generate(spec); err != nil {
		t.Fatalf("生成失败: %v", err)
	}

	for _, path := range []string{
		"internal/admin/model/article.go",
		"internal/admin/handler/article_handler.go",
		"internal/admin/handler/article_handler_test.go",
		"web/apps/admin/src/api/system/article.ts",
		"web/apps/admin/src/views/system/article/index.vue",
		"web/apps/admin/src/views/system/article/modules/form.vue",
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("缺少生成文件 %s: %v", path, err)
		}
	}

	expectContains := map[string]string{
		"internal/admin/module.go":                           "handler.NewArticleHandler(db),",
		"internal/migrate/migrate.go":                        "&adminModel.Article{},",
		"web/apps/admin/src/router/routes/modules/system.ts": "authority: ['system:article:menu']",
	}
	for path, want := range expectContains {
		data, _ := os.ReadFile(path)
		if !strings.Contains(string(data), want) {
			t.Fatalf("%s 应包含 %q，实际:\n%s", path, want, data)
		}
		// 锚点必须保留，后续生成才能继续注册。
		if !strings.Contains(string(data), "// gen:") {
			t.Fatalf("%s 的锚点被删除", path)
		}
	}

	// 重复生成：目标已存在应整体中止，且不重复注册。
	before, _ := os.ReadFile("internal/admin/module.go")
	if err := generate(spec); err == nil {
		t.Fatal("重复生成应返回错误")
	}
	after, _ := os.ReadFile("internal/admin/module.go")
	if string(before) != string(after) {
		t.Fatal("重复生成不应修改注册文件")
	}
}

// TestGenerateMissingAnchor 验证注册文件缺少锚点时报错，而不是静默生成一个不生效的模块。
func TestGenerateMissingAnchor(t *testing.T) {
	setupFakeRepo(t)
	if err := os.WriteFile("internal/migrate/migrate.go", []byte("db.AutoMigrate()\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	spec, _ := newSpec("article", "文章", "title")
	if err := generate(spec); err == nil || !strings.Contains(err.Error(), "gen:models") {
		t.Fatalf("应因缺少 gen:models 锚点报错，实际 %v", err)
	}
}
