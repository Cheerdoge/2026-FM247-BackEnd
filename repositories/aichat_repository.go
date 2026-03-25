package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
)

type AIChatRepository struct {
	redis *redis.Client
}

func NewAIChatRepository(redis *redis.Client) *AIChatRepository {
	return &AIChatRepository{redis: redis}
}

func (r *AIChatRepository) getRedisKey(sessionID uint) string {
	s := strconv.FormatUint(uint64(sessionID), 10)
	return "chat_history:" + s
}

func (r *AIChatRepository) SaveChatHistory(ctx context.Context, sessionID uint, newMessages ...openai.ChatCompletionMessage) error {
	var chatHistory []interface{}
	for _, msg := range newMessages {
		msgJSON, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		chatHistory = append(chatHistory, string(msgJSON))
	}
	rediskey := r.getRedisKey(sessionID)
	err := r.redis.RPush(ctx, rediskey, chatHistory...).Err()
	if err != nil {
		return err
	}
	r.redis.Expire(ctx, rediskey, 7*24*time.Hour)
	return nil
}

func (r *AIChatRepository) GetChatHistory(ctx context.Context, sessionID uint) ([]openai.ChatCompletionMessage, error) {
	rediskey := r.getRedisKey(sessionID)
	values, err := r.redis.LRange(ctx, rediskey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	chatHistory := make([]openai.ChatCompletionMessage, 0, len(values))
	for _, v := range values {
		var message openai.ChatCompletionMessage
		err := json.Unmarshal([]byte(v), &message)
		if err != nil {
			return nil, err
		}
		chatHistory = append(chatHistory, message)
	}
	r.redis.Expire(ctx, rediskey, 7*24*time.Hour)
	return chatHistory, nil
}

// 保持聊天记录不超过50条
func (r *AIChatRepository) TrimChatHistory(ctx context.Context, sessionID uint) error {
	rediskey := r.getRedisKey(sessionID)
	return r.redis.LTrim(ctx, rediskey, -50, -1).Err()
}

// 当ai未回复时删除最新消息，防止污染上下文
func (r *AIChatRepository) PopLatestMessage(ctx context.Context, sessionID uint) error {
	rediskey := r.getRedisKey(sessionID)
	return r.redis.RPop(ctx, rediskey).Err()
}
