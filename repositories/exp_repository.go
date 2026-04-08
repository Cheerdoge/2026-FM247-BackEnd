package repository

import (
	"2026-FM247-BackEnd/models"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ExpRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewExpRepository(db *gorm.DB, redis *redis.Client) *ExpRepository {
	return &ExpRepository{db: db, redis: redis}
}

// 先返回等级，再到返回升级需要的经验值
func CalculateLevel(exp int) (int, int) {
	var levels = []int{0, 0, 30, 60, 120, 180, 300}
	for i := 6; i >= 1; i-- {
		if exp >= levels[i] {
			if i == 6 {
				return 6, 0
			}
			return i, levels[i+1] - exp
		}
	}
	return 0, 0
}

func (r *ExpRepository) UpdateExperience(userID uint, exp int, level int) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{"experience": exp, "level": level}).
		Error
}

func (r *ExpRepository) GetExperienceAndLevelFromDB(userID uint) (int, int, error) {
	var user models.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return 0, 0, result.Error
	}
	return user.Experience, user.Level, nil
}

func (r *ExpRepository) GenerateRedisKey(userID uint) string {
	return fmt.Sprintf("user:state:%d", userID)
}

func (r *ExpRepository) GetExpAndLevelFromRedis(ctx context.Context, userID uint) (int, int, error) {
	rediskey := r.GenerateRedisKey(userID)
	data, err := r.redis.HGetAll(ctx, rediskey).Result()
	if err != nil {
		return 0, 0, err
	}
	if len(data) == 0 {
		exp, level, err := r.GetExperienceAndLevelFromDB(userID)
		if err != nil {
			return 0, 0, err
		}
		pipe := r.redis.Pipeline()
		pipe.HSet(ctx, rediskey, "level", level)
		pipe.HSet(ctx, rediskey, "exp", exp)
		pipe.Expire(ctx, rediskey, 48*time.Hour)
		_, err = pipe.Exec(ctx)
		if err != nil {
			return 0, 0, err
		}
		return exp, level, nil
	}
	r.redis.Expire(ctx, rediskey, 48*time.Hour)
	exp, _ := strconv.Atoi(data["exp"])
	level, _ := strconv.Atoi(data["level"])
	return exp, level, nil
}

func (r *ExpRepository) IncreaseExpAndCheckLevelUp(ctx context.Context, userID uint, exp int) (bool, int, int, error) {
	rediskey := r.GenerateRedisKey(userID)
	if r.redis.Exists(ctx, rediskey).Val() == 0 {
		r.GetExpAndLevelFromRedis(ctx, userID)
	}
	newExp, err := r.redis.HIncrBy(ctx, rediskey, "exp", int64(exp)).Result()
	if err != nil {
		return false, 0, 0, err
	}

	newlevel, expToNextLevel := CalculateLevel(int(newExp))
	oldlevel, _ := strconv.Atoi(r.redis.HGet(ctx, rediskey, "level").Val())
	if newlevel > oldlevel {
		r.redis.HSet(ctx, rediskey, "level", newlevel)
		return true, newlevel, expToNextLevel, nil
	}
	r.redis.Expire(ctx, rediskey, 48*time.Hour)
	return false, newlevel, expToNextLevel, nil
}

func (r *ExpRepository) SyncExpAndLevelToDB(ctx context.Context, userID uint) error {
	rediskey := r.GenerateRedisKey(userID)
	data, err := r.redis.HGetAll(ctx, rediskey).Result()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	exp, _ := strconv.Atoi(data["exp"])
	level, _ := strconv.Atoi(data["level"])
	err = r.UpdateExperience(userID, exp, level)
	if err != nil {
		return err
	}
	r.redis.Expire(ctx, rediskey, 48*time.Hour)
	return nil
}
