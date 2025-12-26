package excel

import (
	"errors"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ParseResult Excel 解析结果。
type ParseResult struct {
	Headers []string
	Rows    [][]string
}

// ParseFromReader 解析 Excel 文件，默认读取第一个 sheet。
//
// 约定：
// - 第一行是表头
// - 从第二行开始是数据行
func ParseFromReader(r io.Reader) (*ParseResult, error) {
	if r == nil {
		return nil, errors.New("文件内容不能为空")
	}

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	sheetName := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("excel 内容为空")
	}

	headers := make([]string, 0, len(rows[0]))
	for _, h := range rows[0] {
		h = strings.TrimSpace(h)
		headers = append(headers, h)
	}

	dataRows := make([][]string, 0)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		// 空行跳过，避免导入无意义数据
		if isEmptyRow(row) {
			continue
		}
		dataRows = append(dataRows, normalizeRow(row, len(headers)))
	}

	return &ParseResult{Headers: headers, Rows: dataRows}, nil
}

// ValidateHeaders 校验表头是否满足期望。
func ValidateHeaders(actual []string, expected []string) error {
	if len(expected) == 0 {
		return errors.New("期望表头不能为空")
	}
	if len(actual) < len(expected) {
		return errors.New("表头列数不匹配")
	}

	for i := 0; i < len(expected); i++ {
		// 允许实际表头后面多列，但前面必须严格匹配
		if strings.TrimSpace(actual[i]) != strings.TrimSpace(expected[i]) {
			return errors.New("表头不匹配")
		}
	}
	return nil
}

func normalizeRow(row []string, n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		if i < len(row) {
			out[i] = strings.TrimSpace(row[i])
			continue
		}
		out[i] = ""
	}
	return out
}

func isEmptyRow(row []string) bool {
	for _, v := range row {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}
