package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	musicService *service.MusicService
}

func NewMusicHandler(musicService *service.MusicService) *MusicHandler {
	return &MusicHandler{
		musicService: musicService,
	}
}

// GetAllMusic 获取所有音乐
// @Router /api/music [get]
func (h *MusicHandler) GetAllMusic(c *gin.Context) {
	musicInfos, err := h.musicService.GetAllMusic()
	if err != nil {
		FailWithMessage(c, "获取音乐列表失败: "+err.Error())
		return
	}
	OkWithData(c, musicInfos)
}

// UploadMusic 上传音乐
// @Router /api/music [post]
func (h *MusicHandler) UploadMusic(c *gin.Context) {
	// 1. 解析表单数据
	var req UploadMusicRequest
	if err := c.ShouldBind(&req); err != nil {
		FailWithMessage(c, "请求参数错误: "+err.Error())
		return
	}

	// 2. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}

	// 3. 获取上传的文件
	file, header, err := c.Request.FormFile("music")
	if err != nil {
		FailWithMessage(c, "请选择头像文件")
		return
	}
	defer file.Close()

	// 4. 上传音乐
	musicURL, err := h.musicService.UploadMusic(c.Request.Context(), claims.UserID, file, header, req.Author, req.Title)
	if err != nil {
		FailWithMessage(c, "上传失败: "+err.Error())
		return
	}

	// 5. 返回音乐URL
	OkWithData(c, gin.H{"music_url": musicURL})
}
