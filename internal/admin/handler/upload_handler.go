package handler

import (
	"bico-admin/internal/core/upload"
	"bico-admin/internal/pkg/crud"
	"bico-admin/internal/pkg/response"
	"errors"

	"github.com/gin-gonic/gin"
)

// UploadHandler 通用上传处理器
//
// 说明：用于富文本编辑器等通用场景上传图片/视频文件。
type UploadHandler struct {
	crud.BaseHandler
	uploader upload.Uploader
}

// NewUploadHandler 创建通用上传处理器
func NewUploadHandler(uploader upload.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

// Upload 上传文件
// @Summary 上传通用文件
// @Description 上传富文本图片或视频文件
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file false "通用文件"
// @Param image formData file false "图片文件"
// @Param video formData file false "视频文件"
// @Param type formData string false "上传类型，image 或 video"
// @Success 200 {object} adminResponse{data=uploadResponse}
// @Router /upload [post]
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
		h.Fail(c, uploadBizError(err))
		return
	}

	response.SuccessWithData(c, gin.H{
		"url": url,
	})
}

// uploadBizError 把上传校验类错误转成业务错误。
//
// 说明：文件过大、类型不支持等由用户操作引起，提示可以直接展示；
// 其余错误（云存储 SDK、磁盘 IO 等）可能包含内部细节，保持原样交给 Fail 记录日志。
func uploadBizError(err error) error {
	for _, target := range []error{
		upload.ErrFileTooLarge,
		upload.ErrInvalidFileType,
		upload.ErrEmptyFile,
		upload.ErrInvalidPath,
	} {
		// 命中任一校验类错误即视为用户可见的业务错误。
		if errors.Is(err, target) {
			return crud.NewBizError(target.Error())
		}
	}
	return err
}
