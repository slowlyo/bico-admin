package handler

import (
	"bico-admin/internal/pkg/crud"
	excelpkg "bico-admin/internal/pkg/excel"
	"bico-admin/internal/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DemoExcelHandler Excel 导入导出示例。
//
// 说明：
// - 模板下载：通用逻辑只生成表头；示例接口会额外写入几行模拟数据
// - 导入：解析并校验表头，返回行数与前几行预览
// - 导出：导出模拟数据并下载
type DemoExcelHandler struct {
	crud.BaseHandler
}

// NewDemoExcelHandler 创建示例 handler。
func NewDemoExcelHandler() *DemoExcelHandler {
	return &DemoExcelHandler{}
}

// DemoExcelHeaders 示例表头定义。
var DemoExcelHeaders = []string{"姓名", "手机号", "年龄", "城市"}

// DownloadTemplate 下载导入模板（示例包含模拟数据）。
func (h *DemoExcelHandler) DownloadTemplate(c *gin.Context) {
	f, err := excelpkg.BuildHeaderTemplate(DemoExcelHeaders)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	_ = excelpkg.AppendRows(f, [][]interface{}{
		{"张三", "13800000000", 20, "上海"},
		{"李四", "13900000000", 28, "北京"},
	})

	filename := "导入模板_示例.xlsx"
	if err := excelpkg.WriteAsAttachment(c, f, filename); err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
}

// Import 导入 Excel（示例）。
func (h *DemoExcelHandler) Import(c *gin.Context) {
	result, err := excelpkg.ParseUploadedAuto(c, "file")
	if err != nil {
		if err == excelpkg.ErrMissingUploadFile {
			h.Error(c, "请上传文件")
			return
		}
		if err == excelpkg.ErrUnsupportedImportType {
			h.Error(c, "不支持的文件类型，请上传 .xlsx/.xlsm/.xltx/.xltm/.csv")
			return
		}
		h.Error(c, err.Error())
		return
	}
	if err := excelpkg.ValidateHeaders(result.Headers, DemoExcelHeaders); err != nil {
		h.Error(c, "导入模板不正确，请先下载模板")
		return
	}

	preview := result.Rows
	if len(preview) > 5 {
		preview = preview[:5]
	}

	h.SuccessWithMessage(c, "导入解析成功", gin.H{
		"total":   len(result.Rows),
		"preview": preview,
	})
}

// Export 导出 Excel（示例）。
func (h *DemoExcelHandler) Export(c *gin.Context) {
	f, err := excelpkg.BuildHeaderTemplate(DemoExcelHeaders)
	if err != nil {
		h.Error(c, err.Error())
		return
	}

	_ = excelpkg.AppendRows(f, [][]interface{}{
		{"王五", "13600000000", 32, "深圳"},
		{"赵六", "13700000000", 25, "杭州"},
		{"钱七", "13500000000", 41, "成都"},
	})

	filename := "导出_示例_" + time.Now().Format("20060102_150405") + ".xlsx"
	if err := excelpkg.WriteAsAttachment(c, f, filename); err != nil {
		response.ErrorWithStatus(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
}
