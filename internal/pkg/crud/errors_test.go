package crud

import (
	"errors"
	"fmt"
	"testing"
)

// TestBizMessage 验证只有业务错误（含被包装的）能提取出可展示的消息。
func TestBizMessage(t *testing.T) {
	sentinel := NewBizError("名称已存在")

	cases := []struct {
		name   string
		err    error
		wantOK bool
		want   string
	}{
		{"直接的业务错误", sentinel, true, "名称已存在"},
		{"被 %w 包装的业务错误", fmt.Errorf("保存失败: %w", sentinel), true, "名称已存在"},
		{"普通错误不暴露消息", errors.New("UNIQUE constraint failed: users.name"), false, ""},
		{"nil 不是业务错误", nil, false, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, ok := bizMessage(tc.err)
			if ok != tc.wantOK || msg != tc.want {
				t.Fatalf("期望 (%q, %v)，实际 (%q, %v)", tc.want, tc.wantOK, msg, ok)
			}
		})
	}
}

// TestBizErrorSentinel 验证业务错误可作为哨兵错误：同一实例 errors.Is 为真，相同文案的不同实例不相等。
func TestBizErrorSentinel(t *testing.T) {
	a := NewBizError("同一文案")
	b := NewBizError("同一文案")
	if !errors.Is(fmt.Errorf("wrap: %w", a), a) {
		t.Fatal("包装后的哨兵错误应能被 errors.Is 识别")
	}
	if errors.Is(a, b) {
		t.Fatal("不同实例按指针比较，不应相等")
	}
	if a.Error() != "同一文案" {
		t.Fatalf("Error() 应返回原文案，实际 %q", a.Error())
	}
}
