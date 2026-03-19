package repository

import (
	"2026-FM247-BackEnd/models"

	"gorm.io/gorm"
)

type AmbientSoundRepository struct {
	db *gorm.DB
}

func NewAmbientSoundRepository(db *gorm.DB) *AmbientSoundRepository {
	return &AmbientSoundRepository{db: db}
}

func (r *AmbientSoundRepository) GetAll() ([]models.AmbientSound, error) {
	var sounds []models.AmbientSound
	result := r.db.Order("created_at desc").Find(&sounds)
	if result.Error != nil {
		return nil, result.Error
	}
	return sounds, nil
}

func (r *AmbientSoundRepository) CreateAmbientSound(name, url string) error {
	sound := models.AmbientSound{
		Name: name,
		URL:  url,
	}
	result := r.db.Create(&sound)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AmbientSoundRepository) DeleteAmbientSound(id uint) error {
	result := r.db.Delete(&models.AmbientSound{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AmbientSoundRepository) GetAmbientSoundByName(name string) (*models.AmbientSound, error) {
	var sound models.AmbientSound
	result := r.db.Where("name = ?", name).First(&sound)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sound, nil
}
