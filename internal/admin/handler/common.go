package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bico-admin/pkg/config"
	"bico-admin/pkg/response"
)

// CommonHandler 通用处理器
type CommonHandler struct {
	cfg *config.Config
}

// NewCommonHandler 创建通用处理器
func NewCommonHandler(cfg *config.Config) *CommonHandler {
	return &CommonHandler{
		cfg: cfg,
	}
}

// UploadResponse 文件上传响应
type UploadResponse struct {
	FileName string `json:"file_name"` // 原始文件名
	FilePath string `json:"file_path"` // 文件相对路径（用于存储）
	FileURL  string `json:"file_url"`  // 文件完整访问URL
	FileSize int64  `json:"file_size"` // 文件大小（字节）
	FileType string `json:"file_type"` // 文件类型
}

// MultiUploadResponse 多文件上传响应
type MultiUploadResponse struct {
	Files []UploadResponse `json:"files"` // 上传成功的文件列表
	Total int              `json:"total"` // 总文件数
}

// Upload 文件上传
// @Summary 文件上传
// @Description 支持单文件或多文件上传，可通过dir参数指定上传目录
// @Tags 通用接口
// @Accept multipart/form-data
// @Produce json
// @Param dir query string false "上传目录，不指定则存储在根目录"
// @Param files formData file true "上传的文件"
// @Success 200 {object} response.ApiResponse{data=MultiUploadResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /admin/upload [post]
func (h *CommonHandler) Upload(c *gin.Context) {
	// 获取上传目录参数
	dir := c.Query("dir")
	if dir == "" {
		dir = c.PostForm("dir") // 也支持通过表单参数传递
	}

	// 解析多文件上传
	form, err := c.MultipartForm()
	if err != nil {
		response.ErrorWithMessage(c, response.CodeBadRequest, "解析上传文件失败: "+err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.ErrorWithMessage(c, response.CodeBadRequest, "请选择要上传的文件")
		return
	}

	// 检查文件数量限制
	maxFiles := h.cfg.Upload.MaxFiles
	if maxFiles > 0 && len(files) > maxFiles {
		response.ErrorWithMessage(c, response.CodeBadRequest, fmt.Sprintf("单次最多只能上传 %d 个文件", maxFiles))
		return
	}

	// 构建上传目录路径
	uploadDir := h.cfg.Upload.UploadDir
	if uploadDir == "" {
		uploadDir = "data/uploads" // 默认目录
	}
	if dir != "" {
		// 清理目录参数，防止路径遍历攻击
		dir = filepath.Clean(dir)
		if strings.Contains(dir, "..") {
			response.ErrorWithMessage(c, response.CodeBadRequest, "无效的目录参数")
			return
		}
		uploadDir = filepath.Join(uploadDir, dir)
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.ErrorWithMessage(c, response.CodeInternalServerError, "创建上传目录失败: "+err.Error())
		return
	}

	var uploadedFiles []UploadResponse
	var failedFiles []string

	// 处理每个上传的文件
	for _, file := range files {
		uploadResp, err := h.saveFile(c, file, uploadDir)
		if err != nil {
			failedFiles = append(failedFiles, fmt.Sprintf("%s: %s", file.Filename, err.Error()))
			continue
		}
		uploadedFiles = append(uploadedFiles, *uploadResp)
	}

	// 如果所有文件都上传失败
	if len(uploadedFiles) == 0 {
		response.ErrorWithMessage(c, response.CodeBadRequest, "所有文件上传失败: "+strings.Join(failedFiles, "; "))
		return
	}

	// 构建响应
	result := MultiUploadResponse{
		Files: uploadedFiles,
		Total: len(uploadedFiles),
	}

	// 如果有部分文件上传失败，在响应中提示
	if len(failedFiles) > 0 {
		message := fmt.Sprintf("成功上传 %d 个文件，%d 个文件上传失败: %s",
			len(uploadedFiles), len(failedFiles), strings.Join(failedFiles, "; "))
		response.SuccessWithMessage(c, message, result)
	} else {
		response.Success(c, result)
	}
}

// saveFile 保存单个文件
func (h *CommonHandler) saveFile(c *gin.Context, file *multipart.FileHeader, uploadDir string) (*UploadResponse, error) {
	// 文件大小验证
	maxFileSize := h.cfg.Upload.GetMaxFileSizeBytes()
	if file.Size > maxFileSize {
		return nil, fmt.Errorf("文件大小超过限制（最大%s）", h.cfg.Upload.MaxFileSize)
	}

	// 文件类型验证
	ext := filepath.Ext(file.Filename)
	if !h.cfg.Upload.IsAllowedType(ext) {
		return nil, fmt.Errorf("不支持的文件类型")
	}

	// 生成唯一文件名
	newFileName := fmt.Sprintf("%s_%s%s",
		time.Now().Format("20060102_150405"),
		uuid.New().String()[:8],
		ext)

	// 构建完整文件路径
	filePath := filepath.Join(uploadDir, newFileName)

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建文件访问路径（转换为可访问的URL路径）
	// 将配置的上传目录转换为URL路径
	accessPath := strings.ReplaceAll(filePath, "\\", "/")
	uploadDirConfig := h.cfg.Upload.UploadDir
	if uploadDirConfig == "" {
		uploadDirConfig = "data/uploads"
	}
	accessPath = strings.Replace(accessPath, uploadDirConfig, "/uploads", 1)

	// 构建完整的文件访问URL
	var fileURL string
	if h.cfg.Upload.BaseURL != "" {
		// 如果配置了基础URL，使用配置的URL
		fileURL = strings.TrimRight(h.cfg.Upload.BaseURL, "/") + accessPath
	} else {
		// 否则使用请求的域名构建URL
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		fileURL = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, accessPath)
	}

	return &UploadResponse{
		FileName: file.Filename,
		FilePath: accessPath,
		FileURL:  fileURL,
		FileSize: file.Size,
		FileType: h.getFileType(file.Filename),
	}, nil
}

// getFileType 获取文件类型描述
func (h *CommonHandler) getFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	typeMap := map[string]string{
		".jpg":  "图片",
		".jpeg": "图片",
		".png":  "图片",
		".gif":  "图片",
		".bmp":  "图片",
		".webp": "图片",
		".svg":  "图片",
		".pdf":  "PDF文档",
		".doc":  "Word文档",
		".docx": "Word文档",
		".xls":  "Excel表格",
		".xlsx": "Excel表格",
		".ppt":  "PowerPoint演示文稿",
		".pptx": "PowerPoint演示文稿",
		".txt":  "文本文件",
		".csv":  "CSV文件",
		".zip":  "压缩文件",
		".rar":  "压缩文件",
		".7z":   "压缩文件",
		".tar":  "压缩文件",
		".gz":   "压缩文件",
		".mp3":  "音频文件",
		".mp4":  "视频文件",
		".avi":  "视频文件",
		".mov":  "视频文件",
		".wmv":  "视频文件",
		".flv":  "视频文件",
		".wav":  "音频文件",
	}
	if fileType, ok := typeMap[ext]; ok {
		return fileType
	}
	return "其他文件"
}
