package excel

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// WriteAsAttachment 将 Excel 写入响应并作为附件下载。
func WriteAsAttachment(c *gin.Context, f *excelize.File, filename string) error {
	if filename == "" {
		filename = "export.xlsx"
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		filename = filename + ".xlsx"
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}

	escaped := url.PathEscape(filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+escaped)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	return nil
}
