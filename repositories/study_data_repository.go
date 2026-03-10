package repository

import (
	"2026-FM247-BackEnd/models"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StudyDataRepository struct {
	db    *gorm.DB
	redis *redis.Client
	ctx   context.Context
}

func NewStudyDataRepository(db *gorm.DB, redis *redis.Client) *StudyDataRepository {
	return &StudyDataRepository{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
}

// 生成redis的每日key
func (r *StudyDataRepository) GenerateDailyKey(userID uint, date time.Time) string {
	return fmt.Sprintf("user:%d:studydata:date:%s", userID, date.Format("2006-01-02"))
}

// 生成redis的每月key
func (r *StudyDataRepository) GenerateMonthlyKey(userID uint, date time.Time) string {
	return fmt.Sprintf("user:%d:studydata:month:%s", userID, date.Format("2006-01"))
}

//redis中仅包含每日学习时长数据，每月将在mysql中聚合
//数据结构： key：user:{userID}:studydata:date:{YYYY-MM-DD}
// 			study_time：学习时长 单位分钟		tomatoes：番茄钟数量

// 增加每日和每月学习时长
func (r *StudyDataRepository) IncrementDailyStudyTime(userID uint, date time.Time, studyTime int) error {
	dailykey := r.GenerateDailyKey(userID, date)
	monthlykey := r.GenerateMonthlyKey(userID, date)

	// 策略：操作前先检查Key是否存在，若不存在尝试从DB回捞数据进行预热
	if r.redis.Exists(r.ctx, dailykey).Val() == 0 {
		r.GetDailyStudyData(userID, date)
	}
	if r.redis.Exists(r.ctx, monthlykey).Val() == 0 {
		r.GetMonthlyStudyData(userID, date)
	}

	pipe := r.redis.Pipeline()
	pipe.HIncrBy(r.ctx, dailykey, "study_time", int64(studyTime))
	pipe.HIncrBy(r.ctx, monthlykey, "study_time", int64(studyTime))

	//设置并只设置一次过期时间(NX)
	pipe.ExpireNX(r.ctx, dailykey, 26*time.Hour)
	pipe.ExpireNX(r.ctx, monthlykey, 35*24*time.Hour)

	_, err := pipe.Exec(r.ctx)
	return err
}

// 增加每日和每月番茄钟次数
func (r *StudyDataRepository) IncrementDailyTomatoes(userID uint, date time.Time, tomatoes int) error {
	dailykey := r.GenerateDailyKey(userID, date)
	monthlykey := r.GenerateMonthlyKey(userID, date)

	// 关键修复：同上，防止番茄钟数据重置覆盖
	if r.redis.Exists(r.ctx, dailykey).Val() == 0 {
		r.GetDailyStudyData(userID, date)
	}
	if r.redis.Exists(r.ctx, monthlykey).Val() == 0 {
		r.GetMonthlyStudyData(userID, date)
	}

	pipe := r.redis.Pipeline()
	pipe.HIncrBy(r.ctx, dailykey, "tomatoes", int64(tomatoes))
	pipe.HIncrBy(r.ctx, monthlykey, "tomatoes", int64(tomatoes))

	//其实可以不写来着，以防万一
	pipe.ExpireNX(r.ctx, dailykey, 26*time.Hour)
	pipe.ExpireNX(r.ctx, monthlykey, 35*24*time.Hour)

	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

// 同步数据到mysql
func (r *StudyDataRepository) SyncDailyDataToMySQL(userID uint, date time.Time, addStudyTime int, addTomatoes int) error {
	dailykey := r.GenerateDailyKey(userID, date)
	monthlykey := r.GenerateMonthlyKey(userID, date)

	daydata, err := r.redis.HGetAll(r.ctx, dailykey).Result()
	if err != nil {
		return err
	}
	if len(daydata) == 0 {
		return nil
	}
	monthdata, err := r.redis.HGetAll(r.ctx, monthlykey).Result()
	if err != nil {
		return err
	}

	//每日的日期格式xxxx-xx-xx-00-00-00
	daystudytime, _ := strconv.Atoi(daydata["study_time"])
	daytomatoes, _ := strconv.Atoi(daydata["tomatoes"])
	dailyData := models.DailyStudyData{
		UserID:    userID,
		Date:      time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()),
		StudyTime: daystudytime,
		Tomatoes:  daytomatoes,
	}
	//每月的日期格式xxxx-xx-01-00-00-00
	monthstudytime, _ := strconv.Atoi(monthdata["study_time"])
	monthtomatoes, _ := strconv.Atoi(monthdata["tomatoes"])
	monthlyData := models.MonthlyStudyData{
		UserID:    userID,
		Month:     time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location()),
		StudyTime: monthstudytime,
		Tomatoes:  monthtomatoes,
	}

	totalData := models.TotalStudyData{
		UserID:    userID,
		StudyTime: addStudyTime,
		Tomatoes:  addTomatoes,
	}

	tx := r.db.Begin()
	//更新每日数据
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"study_time", "tomatoes", "updated_at"}),
	}).Create(&dailyData).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//更新每月数据
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "month"}},
		DoUpdates: clause.AssignmentColumns([]string{"study_time", "tomatoes", "updated_at"}),
	}).Create(&monthlyData).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//更新总数据
	err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"study_time": gorm.Expr("study_time + ?", addStudyTime),
			"tomatoes":   gorm.Expr("tomatoes + ?", addTomatoes),
			"updated_at": time.Now(),
		}),
	}).Create(&totalData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// 事务提交成功后，删除总数据缓存以保证下次读取最新数据
	totalkey := fmt.Sprintf("user:%d:studydata:total", userID)
	r.redis.Del(r.ctx, totalkey)

	return nil
}

//以下是查询数据以及总结报告的数据汇总

//待办：优化过期逻辑，若数据在redis中存在，则延长过期时间

// 查询每日学习数据
// 考虑先从redis中查，没有再从mysql中查，同时将该数据同步至redis，设置过期时间25小时
func (r *StudyDataRepository) GetDailyStudyData(userID uint, date time.Time) (*models.DailyStudyData, error, bool) {
	dailykey := r.GenerateDailyKey(userID, date)

	data, err := r.redis.HGetAll(r.ctx, dailykey).Result()
	if err != nil {
		return nil, err, false
	}

	if len(data) == 0 {
		//查询mysql
		var dailyData models.DailyStudyData
		result := r.db.Where("user_id = ? AND date = ?", userID, time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())).First(&dailyData)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("还未开始记录学习数据"), true
		}
		if result.Error != nil {
			return nil, result.Error, false
		}
		//同步至redis
		pipe := r.redis.Pipeline()
		pipe.HSet(r.ctx, dailykey, "study_time", dailyData.StudyTime)
		pipe.HSet(r.ctx, dailykey, "tomatoes", dailyData.Tomatoes)
		pipe.Expire(r.ctx, dailykey, 26*time.Hour)
		_, err = pipe.Exec(r.ctx)
		if err != nil {
			return nil, err, false
		}
		return &dailyData, nil, false
	}

	studyTime, _ := strconv.Atoi(data["study_time"])
	tomatoes, _ := strconv.Atoi(data["tomatoes"])
	dailyData := &models.DailyStudyData{
		UserID:    userID,
		Date:      time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()),
		StudyTime: studyTime,
		Tomatoes:  tomatoes,
	}

	return dailyData, nil, false
}

// 查询每月学习数据
// 思路同上
func (r *StudyDataRepository) GetMonthlyStudyData(userID uint, date time.Time) (*models.MonthlyStudyData, error, bool) {
	monthlykey := r.GenerateMonthlyKey(userID, date)
	data, err := r.redis.HGetAll(r.ctx, monthlykey).Result()
	if err != nil {
		return nil, err, false
	}
	if len(data) == 0 {
		//查询mysql
		var monthlyData models.MonthlyStudyData
		result := r.db.Where("user_id = ? AND month = ?", userID, time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())).First(&monthlyData)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("还未开始记录学习数据"), true
		}
		if result.Error != nil {
			return nil, result.Error, false
		}
		//同步至redis
		pipe := r.redis.Pipeline()
		pipe.HSet(r.ctx, monthlykey, "study_time", monthlyData.StudyTime)
		pipe.HSet(r.ctx, monthlykey, "tomatoes", monthlyData.Tomatoes)
		pipe.Expire(r.ctx, monthlykey, 35*24*time.Hour)
		_, err = pipe.Exec(r.ctx)
		if err != nil {
			return nil, err, false
		}
		return &monthlyData, nil, false
	}
	studyTime, _ := strconv.Atoi(data["study_time"])
	tomatoes, _ := strconv.Atoi(data["tomatoes"])
	monthlyData := &models.MonthlyStudyData{
		UserID:    userID,
		Month:     time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location()),
		StudyTime: studyTime,
		Tomatoes:  tomatoes,
	}
	return monthlyData, nil, false
}

// 查询总学习数据
// 总数据仅存于mysql中, 但为了性能优化, 查询后会缓存至redis, 设置过期时间24小时
func (r *StudyDataRepository) GetTotalStudyData(userID uint) (*models.TotalStudyData, error, bool) {
	totalkey := fmt.Sprintf("user:%d:studydata:total", userID)
	data, err := r.redis.HGetAll(r.ctx, totalkey).Result()
	if err != nil {
		return nil, err, false
	}
	if len(data) == 0 {
		//查询mysql
		var totalData models.TotalStudyData
		result := r.db.Where("user_id = ?", userID).First(&totalData)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("还未开始记录学习数据"), true
		}
		if result.Error != nil {
			return nil, result.Error, false
		}
		//同步至redis
		pipe := r.redis.Pipeline()
		pipe.HSet(r.ctx, totalkey, "study_time", totalData.StudyTime)
		pipe.HSet(r.ctx, totalkey, "tomatoes", totalData.Tomatoes)
		pipe.ExpireNX(r.ctx, totalkey, 24*time.Hour)
		_, err = pipe.Exec(r.ctx)
		if err != nil {
			return nil, err, false
		}
		return &totalData, nil, false
	}
	studyTime, _ := strconv.Atoi(data["study_time"])
	tomatoes, _ := strconv.Atoi(data["tomatoes"])
	totalData := &models.TotalStudyData{
		UserID:    userID,
		StudyTime: studyTime,
		Tomatoes:  tomatoes,
	}
	return totalData, nil, false
}

// 查询某段时间内的学习数据，进而生成总结报告同时返回具体每日数据，但计算总和交给上层调用者
// 一次数据库范围查询 + Redis获取今日最新数据修正
func (r *StudyDataRepository) GetStudyDataSummary(userID uint, startDate, endDate time.Time) ([]models.DailyStudyData, error) {
	//一次性从MySQL中获取该时间段内的所有每日数据
	var dailyDataList []models.DailyStudyData
	result := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Order("date ASC").Find(&dailyDataList)
	if result.Error != nil {
		return nil, result.Error
	}

	// 构建一个临时Map方便后续由于Redis数据修正时快速查找
	dataMap := make(map[string]*models.DailyStudyData)
	for i := range dailyDataList {
		// 使用指针，这样修改Map中的值会直接影响切片中的元素
		dataMap[dailyDataList[i].Date.Format("2006-01-02")] = &dailyDataList[i]
	}

	// 修正"今天"的数据 (MySQL中的今天可能不是最新的，以Redis为准)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 如果查询范围包含"今天"
	if (today.Equal(startDate) || today.After(startDate)) && (today.Equal(endDate) || today.Before(endDate)) {
		dailyKey := r.GenerateDailyKey(userID, today)
		redisData, err := r.redis.HGetAll(r.ctx, dailyKey).Result()

		if err == nil && len(redisData) > 0 {
			// 解析Redis最新数据
			realTimeStudy, _ := strconv.Atoi(redisData["study_time"])
			realTimeTomato, _ := strconv.Atoi(redisData["tomatoes"])
			todayStr := today.Format("2006-01-02")

			record, exists := dataMap[todayStr]
			if exists {
				// 情况A: 数据库有记录，但Redis是最新的 -> 覆盖
				record.StudyTime = realTimeStudy
				record.Tomatoes = realTimeTomato
			} else {
				// 情况B: 数据库没记录(今天第一次还没同步)，但Redis有 -> 追加
				todayRecord := models.DailyStudyData{
					UserID:    userID,
					Date:      today,
					StudyTime: realTimeStudy,
					Tomatoes:  realTimeTomato,
				}
				dailyDataList = append(dailyDataList, todayRecord)
			}
		}
	}

	return dailyDataList, nil
}
