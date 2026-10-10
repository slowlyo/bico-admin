package pagination

import "testing"

// TestGetOrderBy 验证排序白名单：命中才拼接，方向只认 ascend，其余一律回退为空或降序。
func TestGetOrderBy(t *testing.T) {
	sortable := []string{"id", "name"}
	cases := []struct {
		name  string
		field string
		order string
		want  string
	}{
		{"白名单字段升序", "name", "ascend", "name ASC"},
		{"白名单字段降序", "name", "descend", "name DESC"},
		{"方向缺省按降序", "id", "", "id DESC"},
		{"未知方向按降序", "id", "asc", "id DESC"},
		{"字段不在白名单", "password", "ascend", ""},
		{"注入片段被忽略", "id; DROP TABLE users", "ascend", ""},
		{"字段为空", "", "ascend", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Pagination{SortField: tc.field, SortOrder: tc.order}
			if got := p.GetOrderBy(sortable); got != tc.want {
				t.Fatalf("期望 %q，实际 %q", tc.want, got)
			}
		})
	}

	// 白名单为空时任何字段都不允许排序。
	p := &Pagination{SortField: "id", SortOrder: "ascend"}
	if got := p.GetOrderBy(nil); got != "" {
		t.Fatalf("空白名单应忽略排序，实际 %q", got)
	}
}

// TestPageBounds 验证页码和每页数量的默认值与上限，以及偏移量计算。
func TestPageBounds(t *testing.T) {
	cases := []struct {
		name       string
		page, size int
		wantPage   int
		wantSize   int
		wantOffset int
	}{
		{"正常值", 3, 20, 3, 20, 40},
		{"页码小于 1 回退默认", 0, 20, 1, 20, 0},
		{"每页数量小于 1 回退默认", 2, 0, 2, DefaultPageSize, DefaultPageSize},
		{"每页数量超过上限被截断", 2, 1000, 2, MaxPageSize, MaxPageSize},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Pagination{Page: tc.page, PageSize: tc.size}
			if got := p.GetPage(); got != tc.wantPage {
				t.Fatalf("页码期望 %d，实际 %d", tc.wantPage, got)
			}
			if got := p.GetPageSize(); got != tc.wantSize {
				t.Fatalf("每页数量期望 %d，实际 %d", tc.wantSize, got)
			}
			if got := p.GetOffset(); got != tc.wantOffset {
				t.Fatalf("偏移量期望 %d，实际 %d", tc.wantOffset, got)
			}
		})
	}
}
