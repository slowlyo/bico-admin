package excel

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedImportType = errors.New("不支持的导入文件类型")
)

// ParseAutoFromReader 根据文件名扩展名自动解析。
//
// 支持：
// - xlsx/xlsm/xltx/xltm：按 Excel 解析
// - csv：按 CSV 解析
//
// 不支持：
// - xls：excelize 不支持旧版 xls
func ParseAutoFromReader(r io.Reader, filename string) (*ParseResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// 没扩展名时默认按 Excel 尝试解析
		return ParseFromReader(r)
	}

	switch ext {
	case ".xlsx", ".xlsm", ".xltx", ".xltm":
		return ParseFromReader(r)
	case ".csv":
		return ParseCSVFromReader(r)
	case ".xls":
		return nil, errors.New("暂不支持 .xls 格式，请另存为 .xlsx 或 .csv 后再导入")
	default:
		return nil, ErrUnsupportedImportType
	}
}
