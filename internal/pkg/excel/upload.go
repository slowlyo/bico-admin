package excel

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrMissingUploadFile = errors.New("请上传文件")
)

// GetUploadedFile 从表单中获取上传文件。
func GetUploadedFile(c *gin.Context, field string) (*multipart.FileHeader, error) {
	if c == nil {
		return nil, errors.New("context 不能为空")
	}
	if strings.TrimSpace(field) == "" {
		field = "file"
	}

	file, err := c.FormFile(field)
	if err != nil {
		return nil, ErrMissingUploadFile
	}
	return file, nil
}

// ValidateImportFilename 校验导入文件名对应的扩展名是否支持。
func ValidateImportFilename(filename string) error {
	ext := strings.ToLower(strings.TrimSpace(filename))
	if ext == "" {
		return ErrUnsupportedImportType
	}
	// 这里只做扩展名层面的快速校验，具体解析仍由 ParseAutoFromReader 兜底
	switch {
	case strings.HasSuffix(ext, ".xlsx"),
		strings.HasSuffix(ext, ".xlsm"),
		strings.HasSuffix(ext, ".xltx"),
		strings.HasSuffix(ext, ".xltm"),
		strings.HasSuffix(ext, ".csv"):
		return nil
	case strings.HasSuffix(ext, ".xls"):
		return errors.New("暂不支持 .xls 格式，请另存为 .xlsx 或 .csv 后再导入")
	default:
		return ErrUnsupportedImportType
	}
}

// OpenUploadedFile 打开上传文件并返回可读流（调用方负责 close）。
func OpenUploadedFile(file *multipart.FileHeader) (multipart.File, error) {
	if file == nil {
		return nil, ErrMissingUploadFile
	}
	return file.Open()
}

// ParseUploadedAuto 从上传文件中自动解析（支持 xlsx/xlsm/xltx/xltm/csv）。
func ParseUploadedAuto(c *gin.Context, field string) (*ParseResult, error) {
	file, err := GetUploadedFile(c, field)
	if err != nil {
		return nil, err
	}
	if err := ValidateImportFilename(file.Filename); err != nil {
		return nil, err
	}

	src, err := OpenUploadedFile(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = src.Close()
	}()

	return ParseAutoFromReader(src, file.Filename)
}
