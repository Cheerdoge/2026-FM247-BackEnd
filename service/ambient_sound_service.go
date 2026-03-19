package service

import (
	"2026-FM247-BackEnd/models"
	"2026-FM247-BackEnd/storage"
	"context"
	"errors"
	"mime/multipart"
)

type AmbientSoundRepository interface {
	GetAll() ([]models.AmbientSound, error)
	CreateAmbientSound(name, url string) error
	DeleteAmbientSound(id uint) error
	GetAmbientSoundByName(name string) (*models.AmbientSound, error)
}

type AmbientSoundService struct {
	repo    AmbientSoundRepository
	storage storage.Storage
}

func NewAmbientSoundService(repo AmbientSoundRepository, storage storage.Storage) *AmbientSoundService {
	return &AmbientSoundService{
		repo:    repo,
		storage: storage,
	}
}

func (s *AmbientSoundService) GetAllAmbientSounds() ([]AmbientSoundInfo, error) {
	sounds, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var result []AmbientSoundInfo
	for _, sound := range sounds {
		var info AmbientSoundInfo
		info.Name = sound.Name
		url, err := s.storage.GetURL(sound.URL)
		if err != nil {
			return nil, err
		}
		info.URL = url
		result = append(result, info)
	}
	return result, nil
}

func (s *AmbientSoundService) CreateAmbientSound(ctx context.Context, file multipart.File,
	fileHeader *multipart.FileHeader, name string) (string, error) {

	_, err := s.repo.GetAmbientSoundByName(name)
	if err == nil {
		return "", errors.New("环境音名称已存在")
	}
	uploadpath := "ambient_sounds/" + fileHeader.Filename
	path, err := s.storage.Upload(ctx, uploadpath, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	err = s.repo.CreateAmbientSound(name, path)
	if err != nil {
		return "", err
	}

	fullURL, err := s.storage.GetURL(path)
	if err != nil {
		return "", err
	}

	return fullURL, nil
}

func (s *AmbientSoundService) DeleteAmbientSound(name string) error {
	sound, err := s.repo.GetAmbientSoundByName(name)
	if err != nil {
		return err
	}
	//暂时只删除了数据库中的记录
	err = s.repo.DeleteAmbientSound(sound.ID)
	if err != nil {
		return err
	}
	return nil
}
