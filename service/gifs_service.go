package service

import (
	"2026-FM247-BackEnd/models"
	"2026-FM247-BackEnd/storage"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type GifsRepository interface {
	GetGifURLByID(id uint) (string, error)
	GetGifs() ([]models.Gif, error)
	CreateGif(gif *models.Gif) error
}

type GifsService struct {
	storage storage.Storage
	repo    GifsRepository
}

func NewGifsService(storage storage.Storage, gifsRepo GifsRepository) *GifsService {
	return &GifsService{storage: storage, repo: gifsRepo}
}

func (s *GifsService) GetGifURLByID(id uint) (string, error) {
	url, err := s.repo.GetGifURLByID(id)
	if err != nil {
		return "", err
	}
	gifurl, err := s.storage.GetURL(url)
	if err != nil {
		return "", err
	}
	return gifurl, nil
}

func (s *GifsService) GetGifs() ([]GifInfo, error) {
	gifs, err := s.repo.GetGifs()
	if err != nil {
		return nil, err
	}
	var gifInfos []GifInfo
	for _, gif := range gifs {
		url, err := s.storage.GetURL(gif.URL)
		if err != nil {
			return nil, err
		}
		gifInfos = append(gifInfos, GifInfo{
			ID:   gif.ID,
			Name: gif.Name,
			URL:  url,
		})
	}
	return gifInfos, nil
}

func (s *GifsService) CreateGif(ctx context.Context, name string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	pathname := filepath.Base(name)
	path := fmt.Sprintf("gifs/%s", pathname)
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/gif") {
		return "", errors.New("只允许上传gif文件")
	}
	url, err := s.storage.Upload(ctx, path, file, fileHeader.Size, contentType)
	if err != nil {
		return "", err
	}
	gif := &models.Gif{
		Name: name,
		URL:  url,
	}
	gifurl, err := s.storage.GetURL(url)
	if err != nil {
		return "", err
	}
	err = s.repo.CreateGif(gif)
	if err != nil {
		return "", err
	}
	return gifurl, err
}
