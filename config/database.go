package config

import (
	"fmt"
	"log"

	"2026-FM247-BackEnd/gormlogger"
	"2026-FM247-BackEnd/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormlogger.NewStdLogger(logger.Info)})
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("Database connection established")

	// 自动迁移
	db.AutoMigrate(
		&models.User{},
		&models.TotalStudyData{},
		&models.DailyStudyData{},
		&models.MonthlyStudyData{},
		&models.Todo{},
		&models.Note{},
		&models.TokenBlacklist{},
		&models.Music{},
		&models.AmbientSound{},
	)
	log.Println("Database migrated successfully")
	return db, nil
}

func CloseDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}
