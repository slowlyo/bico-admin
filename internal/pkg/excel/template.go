package excel

import (
	"errors"

	"github.com/xuri/excelize/v2"
)

// BuildHeaderTemplate 构建只包含表头的 Excel 模板。
func BuildHeaderTemplate(headers []string) (*excelize.File, error) {
	if len(headers) == 0 {
		return nil, errors.New("表头不能为空")
	}

	f := excelize.NewFile()
	sheetName := f.GetSheetName(f.GetActiveSheetIndex())

	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, err
		}
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, err
		}
	}

	_ = f.SetRowHeight(sheetName, 1, 20)
	_ = f.SetColWidth(sheetName, "A", "Z", 18)

	return f, nil
}

// AppendRows 追加数据行（用于示例或业务导出）。
func AppendRows(f *excelize.File, rows [][]interface{}) error {
	if f == nil {
		return errors.New("excel 文件不能为空")
	}
	if len(rows) == 0 {
		return nil
	}

	sheetName := f.GetSheetName(f.GetActiveSheetIndex())
	for r, row := range rows {
		for c, v := range row {
			cell, err := excelize.CoordinatesToCellName(c+1, r+2)
			if err != nil {
				return err
			}
			if err := f.SetCellValue(sheetName, cell, v); err != nil {
				return err
			}
		}
	}

	return nil
}
