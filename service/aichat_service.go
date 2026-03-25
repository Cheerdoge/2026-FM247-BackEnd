package service

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

const MaxChatHistoryMessages = 20

type AIChatRepository interface {
	SaveChatHistory(ctx context.Context, sessionID uint, newMessages ...openai.ChatCompletionMessage) error
	GetChatHistory(ctx context.Context, sessionID uint) ([]openai.ChatCompletionMessage, error)
	TrimChatHistory(ctx context.Context, sessionID uint) error
	PopLatestMessage(ctx context.Context, sessionID uint) error
}

type AIChatService struct {
	repo AIChatRepository
	ai   *openai.Client
}

func NewAIChatService(repo AIChatRepository, ai *openai.Client) *AIChatService {
	return &AIChatService{repo: repo, ai: ai}
}

func (s *AIChatService) Chat(ctx context.Context, userID uint, content string) (string, error) {
	userMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}

	err := s.repo.SaveChatHistory(ctx, userID, userMsg)
	if err != nil {
		return "", fmt.Errorf("存入聊天记录失败: %w", err)
	}

	aiRequestSuccess := false
	defer func() {
		if !aiRequestSuccess {
			s.repo.PopLatestMessage(ctx, userID)
		}
	}()

	chatHistory, err := s.repo.GetChatHistory(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("获取聊天记录失败: %w", err)
	}

	finalMessages := make([]openai.ChatCompletionMessage, 0, MaxChatHistoryMessages+1)
	systemPrompt := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "你是一个线上自习室电台的主持人，名字叫Monica，性格温和、善解人意、独立、沉稳，负责与用户进行友好、温暖的对话，提供学习建议和心理支持。请根据用户的提问，结合自习室的氛围，给出有帮助的回答。你可以分享一些学习方法、时间管理技巧，或者只是陪伴用户聊天，缓解他们的压力。请保持语气亲切、鼓励和理解，让用户感受到温暖和支持。",
	}
	finalMessages = append(finalMessages, systemPrompt)
	startIndex := 0
	if len(chatHistory) > MaxChatHistoryMessages {
		startIndex = len(chatHistory) - MaxChatHistoryMessages
	}
	finalMessages = append(finalMessages, chatHistory[startIndex:]...)

	resp, err := s.ai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       "qwen-plus",
		Messages:    finalMessages,
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("调用AI接口失败: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI接口返回空结果: request_id=%s, model=%s", resp.ID, resp.Model)
	}
	aiMsg := resp.Choices[0].Message
	err = s.repo.SaveChatHistory(ctx, userID, aiMsg)
	if err != nil {
		return "", fmt.Errorf("保存AI回复失败: %w", err)
	}
	aiRequestSuccess = true
	return aiMsg.Content, nil
}

func (s *AIChatService) GetChatHistory(ctx context.Context, userID uint) ([]openai.ChatCompletionMessage, error) {
	chatHistory, err := s.repo.GetChatHistory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天记录失败: %w", err)
	}
	return chatHistory, nil
}
