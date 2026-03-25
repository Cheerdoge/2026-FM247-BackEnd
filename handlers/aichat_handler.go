package handler

import (
	"2026-FM247-BackEnd/utils"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type AIChatService interface {
	Chat(ctx context.Context, userID uint, content string) (string, error)
	GetChatHistory(ctx context.Context, userID uint) ([]openai.ChatCompletionMessage, error)
}

type AIChatHandler struct {
	service AIChatService
}

func NewAIChatHandler(service AIChatService) *AIChatHandler {
	return &AIChatHandler{service: service}
}

func (h *AIChatHandler) Chat(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		FailWithMessage(c, "请求参数错误")
		return
	}

	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	response, err := h.service.Chat(c.Request.Context(), claims.UserID, req.Content)
	if err != nil {
		FailWithMessage(c, "聊天请求失败")
		fmt.Printf("Chat error: %v\n", err)
		return
	}
	OkWithData(c, response)
}

func (h *AIChatHandler) GetChatHistory(c *gin.Context) {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}
	history, err := h.service.GetChatHistory(c.Request.Context(), claims.UserID)
	if err != nil {
		FailWithMessage(c, "获取聊天记录失败")
		return
	}
	OkWithData(c, history)
}
