package repository

import (
	"2026-FM247-BackEnd/models"

	"gorm.io/gorm"
)

type GifsRepository struct {
	db *gorm.DB
}

func NewGifsRepository(db *gorm.DB) *GifsRepository {
	return &GifsRepository{db: db}
}

func (r *GifsRepository) GetGifURLByID(id uint) (string, error) {
	var gif models.Gif
	err := r.db.First(&gif, id).Error
	if err != nil {
		return "", err
	}
	return gif.URL, nil
}

func (r *GifsRepository) GetGifs() ([]models.Gif, error) {
	var gifs []models.Gif
	err := r.db.Find(&gifs).Error
	if err != nil {
		return nil, err
	}
	return gifs, nil
}

func (r *GifsRepository) CreateGif(gif *models.Gif) error {
	return r.db.FirstOrCreate(gif, models.Gif{Name: gif.Name}).Error
}
