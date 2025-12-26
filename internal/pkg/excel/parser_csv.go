package excel

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

// ParseCSVFromReader 解析 CSV 文件。
//
// 约定：
// - 第一行是表头
// - 从第二行开始是数据行
func ParseCSVFromReader(r io.Reader) (*ParseResult, error) {
	if r == nil {
		return nil, errors.New("文件内容不能为空")
	}

	reader := csv.NewReader(bufio.NewReader(r))
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, errors.New("csv 内容为空")
	}

	headers := make([]string, 0, len(records[0]))
	for i, h := range records[0] {
		// 处理 UTF-8 BOM
		if i == 0 {
			h = strings.TrimPrefix(h, "\ufeff")
		}
		headers = append(headers, strings.TrimSpace(h))
	}

	dataRows := make([][]string, 0)
	for i := 1; i < len(records); i++ {
		row := records[i]
		if isEmptyRow(row) {
			continue
		}
		dataRows = append(dataRows, normalizeRow(row, len(headers)))
	}

	return &ParseResult{Headers: headers, Rows: dataRows}, nil
}
