package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type AmbientSoundService interface {
	GetAllAmbientSounds() ([]service.AmbientSoundInfo, error)
	CreateAmbientSound(ctx context.Context, file multipart.File,
		fileHeader *multipart.FileHeader, title string) (string, error)
	DeleteAmbientSound(name string) error
}

type AmbientSoundHandler struct {
	service AmbientSoundService
}

func NewAmbientSoundHandler(service AmbientSoundService) *AmbientSoundHandler {
	return &AmbientSoundHandler{
		service: service,
	}
}

// GetAllAmbientSounds 获取所有环境音
// @Router /api/ambient-sounds [get]
func (h *AmbientSoundHandler) GetAllAmbientSounds(c *gin.Context) {
	sounds, err := h.service.GetAllAmbientSounds()
	if err != nil {
		FailWithMessage(c, "获取环境音列表失败: "+err.Error())
		return
	}
	OkWithData(c, sounds)
}

// CreateAmbientSound 创建环境音
// @Router /api/ambient-sounds [post]
func (h *AmbientSoundHandler) CreateAmbientSound(c *gin.Context) {
	// 1. 解析表单数据
	var req CreateAmbientSoundRequest
	if err := c.ShouldBind(&req); err != nil {
		FailWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 2. 验证登录
	_, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}

	// 2. 获取上传的文件
	file, header, err := c.Request.FormFile("sound")
	if err != nil {
		FailWithMessage(c, "请选择环境音文件")
		return
	}
	defer file.Close()

	// 3. 上传环境音
	soundURL, err := h.service.CreateAmbientSound(c.Request.Context(), file, header, req.Name)
	if err != nil {
		FailWithMessage(c, "上传失败: "+err.Error())
		return
	}
	OkWithData(c, gin.H{"url": soundURL})
}

// DeleteAmbientSound 删除环境音
// @Router /api/ambient-sounds/:name [delete]
func (h *AmbientSoundHandler) DeleteAmbientSound(c *gin.Context) {
	// 1. 解析请求参数
	name := c.Param("name")
	if name == "" {
		FailWithMessage(c, "请提供环境音名称")
		return
	}

	// 2. 验证登录
	_, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}

	// 3. 删除环境音
	err = h.service.DeleteAmbientSound(name)
	if err != nil {
		FailWithMessage(c, "删除失败: "+err.Error())
		return
	}

	OkWithMessage(c, "成功删除")
}
