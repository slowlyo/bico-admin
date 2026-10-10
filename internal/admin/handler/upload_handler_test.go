package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bico-admin/internal/admin/service"
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/crud"

	"github.com/gin-gonic/gin"
)

// failBody 用 BaseHandler.Fail 渲染错误并返回状态码和响应体。
func failBody(err error) (int, string) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	(&crud.BaseHandler{}).Fail(c, err)
	return w.Code, w.Body.String()
}

// TestUploadBizError 验证上传校验类错误对用户可见，其余错误保持原样交给 Fail 隐藏。
func TestUploadBizError(t *testing.T) {
	for _, target := range []error{
		upload.ErrFileTooLarge,
		upload.ErrInvalidFileType,
		upload.ErrEmptyFile,
		upload.ErrInvalidPath,
	} {
		t.Run(target.Error(), func(t *testing.T) {
			// 被包装后的校验错误同样需要识别。
			_, body := failBody(uploadBizError(fmt.Errorf("上传: %w", target)))
			if want := target.Error(); !strings.Contains(body, want) {
				t.Fatalf("响应应包含 %q，实际 %s", want, body)
			}
		})
	}

	// 云存储等系统错误不能把内部细节返回给客户端。
	system := errors.New("oss: AccessKeySecret=abc123 invalid")
	if got := uploadBizError(system); got != system {
		t.Fatalf("非校验类错误应原样返回，实际 %v", got)
	}
	_, body := failBody(uploadBizError(system))
	if strings.Contains(body, "abc123") || !strings.Contains(body, "服务器内部错误") {
		t.Fatalf("系统错误应被隐藏，实际 %s", body)
	}
}

// TestServiceBusinessErrorsVisible 验证登录与账号相关的哨兵错误作为业务错误展示原文案。
func TestServiceBusinessErrorsVisible(t *testing.T) {
	for _, err := range []error{
		service.ErrUserNotFound,
		service.ErrInvalidCredentials,
		service.ErrUserDisabled,
		service.ErrOldPasswordWrong,
		service.ErrLoginLocked,
	} {
		t.Run(err.Error(), func(t *testing.T) {
			code, body := failBody(err)
			if code != http.StatusOK && code != http.StatusBadRequest {
				t.Fatalf("业务错误不应返回 5xx，实际 %d", code)
			}
			if !strings.Contains(body, err.Error()) {
				t.Fatalf("响应应展示 %q，实际 %s", err.Error(), body)
			}
		})
	}
}
