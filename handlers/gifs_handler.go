package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type GifsService interface {
	GetGifURLByID(id uint) (string, error)
	GetGifs() ([]service.GifInfo, error)
	CreateGif(ctx context.Context, name string, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type GifsHandler struct {
	service GifsService
}

func NewGifsHandler(service GifsService) *GifsHandler {
	return &GifsHandler{service: service}
}

// GetGifURLByID 根据ID获取GIF URL
// @Router /gifs/{id} [get]
func (h *GifsHandler) GetGifURLByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.StringToUint(idStr)
	if err != nil {
		FailWithMessage(c, "参数错误")
		return
	}
	url, err := h.service.GetGifURLByID(id)
	if err != nil {
		FailWithMessage(c, "获取失败")
		return
	}
	OkWithData(c, url)
}

// GetGifs 获取所有GIF
// @Router /gifs [get]
func (h *GifsHandler) GetGifs(c *gin.Context) {
	gifs, err := h.service.GetGifs()
	if err != nil {
		FailWithMessage(c, "获取失败")
		return
	}
	OkWithData(c, gifs)
}

// CreateGif 创建新的GIF
// @Router /gifs [post]
func (h *GifsHandler) CreateGif(c *gin.Context) {
	var req CreateGifRequest
	if err := c.ShouldBind(&req); err != nil {
		FailWithMessage(c, "请求参数错误: "+err.Error())
		return
	}
	file, header, err := c.Request.FormFile("gif")
	if err != nil {
		FailWithMessage(c, "请选择gif文件")
		return
	}
	defer file.Close()
	gifURL, err := h.service.CreateGif(c.Request.Context(), req.Name, file, header)
	if err != nil {
		FailWithMessage(c, "上传失败: "+err.Error())
		return
	}
	OkWithData(c, gin.H{"url": gifURL})
}
