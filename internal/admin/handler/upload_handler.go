package handler

import (
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadHandler 通用上传处理器
//
// 说明：用于富文本编辑器等通用场景上传图片/视频文件。
type UploadHandler struct {
	uploader upload.Uploader
}

// NewUploadHandler 创建通用上传处理器
func NewUploadHandler(uploader upload.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

// Upload 上传文件
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		// 兼容部分前端实现可能使用 image/video 作为字段名
		file, err = c.FormFile("image")
		if err != nil {
			// 如果 image 不存在，则再尝试 video 字段
			file, err = c.FormFile("video")
		}
	}
	if err != nil {
		response.BadRequest(c, "请上传文件")
		return
	}

	uploadType := c.PostForm("type")
	if uploadType == "" {
		// 未传 type 时默认按图片处理
		uploadType = "image"
	}

	subPath := "uploads"
	if uploadType == "video" {
		// 上传视频
		subPath = "videos"
	} else {
		// 其他类型统一当做图片
		// 非 video 默认当做图片
		subPath = "images"
	}

	url, err := h.uploader.Upload(file, subPath)
	if err != nil {
		response.ErrorWithCode(c, 400, err.Error())
		return
	}

	response.SuccessWithData(c, gin.H{
		"url": url,
	})
}
