// gen 为标准 CRUD 模块生成前后端样板代码。
//
// 用法：go run ./cmd/gen -name article -title 文章 -fields "title:string:标题,content:text:内容,sort:int:排序"
package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

// Field 描述模块的一个业务字段。
type Field struct {
	Name   string // 下划线列名，如 sort_order
	GoName string // Go 字段名，如 SortOrder
	Type   string // string / text / int / bool
	Label  string // 界面展示名称
}

// GoType 返回字段对应的 Go 类型。
func (f Field) GoType() string {
	switch f.Type {
	case "int":
		return "int"
	case "bool":
		return "bool"
	default:
		return "string"
	}
}

// GormTag 返回字段的 gorm 标签，保证列类型与长度合理。
func (f Field) GormTag() string {
	switch f.Type {
	case "text":
		return "type:text"
	case "int":
		return "not null;default:0"
	case "bool":
		return "not null;default:false"
	default:
		return "size:255;not null"
	}
}

// SampleJSON 返回测试创建请求里该字段的 JSON 值。
func (f Field) SampleJSON() string {
	switch f.Type {
	case "int":
		return "1"
	case "bool":
		return "true"
	default:
		return `"示例"`
	}
}

// UpdatedJSON 返回测试更新请求里该字段的 JSON 值；刻意使用零值（bool 为 false）验证零值也能写入。
func (f Field) UpdatedJSON() string {
	switch f.Type {
	case "int":
		return "2"
	case "bool":
		return "false"
	default:
		return `"已更新"`
	}
}

// UpdatedGo 返回与 UpdatedJSON 对应的 Go 字面量，用于断言数据库结果。
func (f Field) UpdatedGo() string {
	if f.Type == "string" || f.Type == "text" {
		return `"已更新"`
	}
	return f.UpdatedJSON()
}

// IsString 判断字段是否参与列表关键字搜索。
func (f Field) IsString() bool { return f.Type == "string" }

// Spec 是模板渲染使用的模块描述。
type Spec struct {
	Name   string // 下划线单数名，如 article_tag
	Title  string // 中文名称，如 文章标签
	Pascal string // ArticleTag
	Camel  string // articleTag
	Kebab  string // article-tag
	Group  string // 路由分组，如 /article-tags
	Table  string // 表名，如 article_tags
	Fields []Field
}

// First 返回第一个业务字段，测试用它验证更新。
func (s Spec) First() Field { return s.Fields[0] }

// SearchFields 返回参与列表搜索的字符串字段。
func (s Spec) SearchFields() []Field {
	var result []Field
	for _, f := range s.Fields {
		if f.IsString() {
			result = append(result, f)
		}
	}
	return result
}

// SortableColumns 返回允许排序的列名，文本字段不参与排序。
func (s Spec) SortableColumns() []string {
	cols := []string{"id"}
	for _, f := range s.Fields {
		if f.Type != "text" {
			cols = append(cols, f.Name)
		}
	}
	return append(cols, "enabled", "created_at", "updated_at")
}

var (
	namePattern  = regexp.MustCompile(`^[a-z][a-z0-9]*(_[a-z0-9]+)*$`)
	fieldTypes   = map[string]bool{"string": true, "text": true, "int": true, "bool": true}
	reservedCols = map[string]bool{"id": true, "enabled": true, "created_at": true, "updated_at": true}
)

// main 解析参数、渲染模板并把模块注册到装配点。
func main() {
	name := flag.String("name", "", "模块名，下划线单数，如 article_tag")
	title := flag.String("title", "", "中文名称，如 文章标签")
	fields := flag.String("fields", "", "字段列表，格式 name:type[:label]，type 可选 string/text/int/bool，逗号分隔")
	flag.Parse()

	spec, err := newSpec(*name, *title, *fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, "参数错误:", err)
		os.Exit(1)
	}
	if err := generate(spec); err != nil {
		fmt.Fprintln(os.Stderr, "生成失败:", err)
		os.Exit(1)
	}
	fmt.Printf("已生成模块 %s。下一步: make swagger web-types，再执行 make check\n", spec.Name)
}

// newSpec 校验输入并推导各种命名形式。
func newSpec(name, title, fieldsArg string) (Spec, error) {
	if !namePattern.MatchString(name) {
		return Spec{}, fmt.Errorf("-name 必须是下划线小写单数名，如 article_tag")
	}
	if title == "" {
		return Spec{}, fmt.Errorf("-title 不能为空")
	}
	fields, err := parseFields(fieldsArg)
	if err != nil {
		return Spec{}, err
	}
	pascal := toPascal(name)
	return Spec{
		Name:   name,
		Title:  title,
		Pascal: pascal,
		Camel:  strings.ToLower(pascal[:1]) + pascal[1:],
		Kebab:  strings.ReplaceAll(name, "_", "-"),
		Group:  "/" + strings.ReplaceAll(name, "_", "-") + "s",
		Table:  name + "s",
		Fields: fields,
	}, nil
}

// parseFields 解析 name:type[:label] 列表，缺省类型为 string，缺省标签为列名。
func parseFields(arg string) ([]Field, error) {
	if strings.TrimSpace(arg) == "" {
		return nil, fmt.Errorf("-fields 至少需要一个字段")
	}
	seen := map[string]bool{}
	var result []Field
	for _, item := range strings.Split(arg, ",") {
		parts := strings.Split(strings.TrimSpace(item), ":")
		f := Field{Name: parts[0], Type: "string", Label: parts[0]}
		if len(parts) > 1 && parts[1] != "" {
			f.Type = parts[1]
		}
		if len(parts) > 2 && parts[2] != "" {
			f.Label = parts[2]
		}
		if !namePattern.MatchString(f.Name) || reservedCols[f.Name] {
			return nil, fmt.Errorf("字段名 %q 非法或为保留列", f.Name)
		}
		if !fieldTypes[f.Type] {
			return nil, fmt.Errorf("字段 %s 的类型 %q 不支持，可选 string/text/int/bool", f.Name, f.Type)
		}
		if seen[f.Name] {
			return nil, fmt.Errorf("字段 %s 重复", f.Name)
		}
		seen[f.Name] = true
		f.GoName = toPascal(f.Name)
		result = append(result, f)
	}
	return result, nil
}

// toPascal 把下划线命名转成大驼峰。
func toPascal(s string) string {
	var b strings.Builder
	for _, part := range strings.Split(s, "_") {
		if part != "" {
			b.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return b.String()
}

// generate 渲染全部文件并注册到模块装配点；任一目标已存在时整体中止，避免覆盖手写代码。
func generate(spec Spec) error {
	files := map[string]string{
		"templates/model.go.tmpl":        fmt.Sprintf("internal/admin/model/%s.go", spec.Name),
		"templates/handler.go.tmpl":      fmt.Sprintf("internal/admin/handler/%s_handler.go", spec.Name),
		"templates/handler_test.go.tmpl": fmt.Sprintf("internal/admin/handler/%s_handler_test.go", spec.Name),
		"templates/api.ts.tmpl":          fmt.Sprintf("web/apps/admin/src/api/system/%s.ts", spec.Kebab),
		"templates/index.vue.tmpl":       fmt.Sprintf("web/apps/admin/src/views/system/%s/index.vue", spec.Kebab),
		"templates/form.vue.tmpl":        fmt.Sprintf("web/apps/admin/src/views/system/%s/modules/form.vue", spec.Kebab),
	}
	for _, target := range files {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("%s 已存在", target)
		}
	}

	rendered := map[string][]byte{}
	for tmplPath, target := range files {
		content, err := render(tmplPath, spec)
		if err != nil {
			return err
		}
		rendered[target] = content
	}
	for target, content := range rendered {
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return err
		}
	}

	// 注册点缺少锚点时要报错，否则新模块会生成了却不生效。
	registrations := []struct{ file, anchor, line string }{
		{"internal/admin/module.go", "// gen:modules", fmt.Sprintf("handler.New%sHandler(db),", spec.Pascal)},
		{"internal/migrate/migrate.go", "// gen:models", fmt.Sprintf("&adminModel.%s{},", spec.Pascal)},
	}
	for _, r := range registrations {
		if err := insertBeforeAnchor(r.file, r.anchor, r.line); err != nil {
			return err
		}
	}
	routeSnippet, err := render("templates/route.ts.tmpl", spec)
	if err != nil {
		return err
	}
	return insertBeforeAnchor("web/apps/admin/src/router/routes/modules/system.ts", "// gen:routes", strings.TrimRight(string(routeSnippet), "\n"))
}

// render 渲染模板；Go 文件额外经过 gofmt，保证生成结果可直接提交。
func render(tmplPath string, spec Spec) ([]byte, error) {
	// 使用 [[ ]] 作为分隔符，避免与 Vue 模板里的 {{ }} 冲突。
	tmpl, err := template.New(filepath.Base(tmplPath)).Delims("[[", "]]").ParseFS(templates, tmplPath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, spec); err != nil {
		return nil, fmt.Errorf("渲染 %s: %w", tmplPath, err)
	}
	if strings.HasSuffix(tmplPath, ".go.tmpl") {
		return format.Source(buf.Bytes())
	}
	return buf.Bytes(), nil
}

// insertBeforeAnchor 在锚点所在行之前插入内容，并沿用锚点的缩进。
func insertBeforeAnchor(file, anchor, content string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != anchor {
			continue
		}
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		var inserted []string
		for _, l := range strings.Split(content, "\n") {
			inserted = append(inserted, indent+l)
		}
		lines = append(lines[:i], append(inserted, lines[i:]...)...)
		return os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0o644)
	}
	return fmt.Errorf("%s 缺少锚点 %s", file, anchor)
}
