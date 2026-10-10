package handler

import (
	"bico-admin/internal/pkg/crud"
	excelpkg "bico-admin/internal/pkg/excel"
	"strconv"
	"strings"
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

type demoExcelRow struct {
	ID    uint
	Name  string
	Phone string
	Age   int
	City  string
}

var demoExcelRows = []demoExcelRow{
	{ID: 1, Name: "王五", Phone: "13600000000", Age: 32, City: "深圳"},
	{ID: 2, Name: "赵六", Phone: "13700000000", Age: 25, City: "杭州"},
	{ID: 3, Name: "钱七", Phone: "13500000000", Age: 41, City: "成都"},
}

// DownloadTemplate 下载导入模板（示例包含模拟数据）。
// @Summary 下载 Excel 导入模板
// @Description 下载示例 Excel 模板文件
// @Tags 示例
// @Produce application/octet-stream
// @Security BearerAuth
// @Success 200 {file} file
// @Router /demo/excel/template [get]
func (h *DemoExcelHandler) DownloadTemplate(c *gin.Context) {
	f, err := excelpkg.BuildHeaderTemplate(DemoExcelHeaders)
	if err != nil {
		h.Fail(c, err)
		return
	}

	_ = excelpkg.AppendRows(f, [][]interface{}{
		{"张三", "13800000000", 20, "上海"},
		{"李四", "13900000000", 28, "北京"},
	})

	filename := "导入模板_示例.xlsx"
	if err := excelpkg.WriteAsAttachment(c, f, filename); err != nil {
		h.Fail(c, err)
		return
	}
}

// Import 导入 Excel（示例）。
// @Summary 导入 Excel
// @Description 上传并解析示例 Excel 文件
// @Tags 示例
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Excel 或 CSV 文件"
// @Success 200 {object} adminResponse{data=demoExcelImportResponse}
// @Router /demo/excel/import [post]
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
		// 解析失败几乎都是文件内容问题（空文件、格式损坏等），提示面向用户。
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
// @Summary 导出 Excel
// @Description 导出示例 Excel 文件，传 ids 时只导出勾选行
// @Tags 示例
// @Produce application/octet-stream
// @Security BearerAuth
// @Success 200 {file} file
// @Router /demo/excel/export [get]
func (h *DemoExcelHandler) Export(c *gin.Context) {
	f, err := excelpkg.BuildHeaderTemplate(DemoExcelHeaders)
	if err != nil {
		h.Fail(c, err)
		return
	}

	rows := filterDemoExcelRows(parseDemoExcelIDs(c.Query("ids")))
	_ = excelpkg.AppendRows(f, buildDemoExcelExportRows(rows))

	filename := "导出_示例_" + time.Now().Format("20060102_150405") + ".xlsx"
	if err := excelpkg.WriteAsAttachment(c, f, filename); err != nil {
		h.Fail(c, err)
		return
	}
}

// parseDemoExcelIDs 解析导出选中行 ID，非法片段直接忽略。
func parseDemoExcelIDs(value string) []uint {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	ids := make([]uint, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		// 导出勾选行只接受正整数 ID，非法值不参与筛选。
		if err != nil || id == 0 {
			continue
		}
		ids = append(ids, uint(id))
	}
	return crud.UniqueUints(ids)
}

// filterDemoExcelRows 按 ID 过滤示例数据；未传 ID 时导出全部示例数据。
func filterDemoExcelRows(ids []uint) []demoExcelRow {
	if len(ids) == 0 {
		return demoExcelRows
	}

	idSet := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	rows := make([]demoExcelRow, 0, len(ids))
	for _, row := range demoExcelRows {
		// 只导出前端勾选的 ID，保持导出范围和列表选择一致。
		if _, ok := idSet[row.ID]; ok {
			rows = append(rows, row)
		}
	}
	return rows
}

// buildDemoExcelExportRows 转换示例数据为 Excel 行数据。
func buildDemoExcelExportRows(rows []demoExcelRow) [][]interface{} {
	result := make([][]interface{}, 0, len(rows))
	for _, row := range rows {
		result = append(result, []interface{}{row.Name, row.Phone, row.Age, row.City})
	}
	return result
}
