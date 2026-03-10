package repository

import (
	"2026-FM247-BackEnd/models"

	"gorm.io/gorm"
)

type MusicRepository struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) *MusicRepository {
	return &MusicRepository{db: db}
}

func (r *MusicRepository) GetAll() ([]models.Music, error) {
	var musics []models.Music
	result := r.db.Order("created_at desc").Find(&musics)
	return musics, result.Error
}

func (r *MusicRepository) CreateMusic(author, title string, duration int, url string, uploaderID uint) error {
	music := models.Music{
		Author:     author,
		Title:      title,
		Duration:   duration,
		FileURL:    url,
		UploaderID: uploaderID,
	}
	result := r.db.Create(&music)
	return result.Error
}
