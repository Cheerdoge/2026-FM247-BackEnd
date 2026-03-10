package service

import (
	"2026-FM247-BackEnd/models"
	"2026-FM247-BackEnd/storage"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/tcolgate/mp3"
)

type MusicRepository interface {
	GetAll() ([]models.Music, error)
	CreateMusic(author, title string, duration int, url string, uploaderID uint) error
}

type MusicService struct {
	storage   storage.Storage
	musicRepo MusicRepository
}

func NewMusicService(musicRepo MusicRepository, storage storage.Storage) *MusicService {
	return &MusicService{
		storage:   storage,
		musicRepo: musicRepo,
	}
}

// GetAllMusic 获取所有音乐
func (s *MusicService) GetAllMusic() ([]MusicInfo, error) {

	musics, err := s.musicRepo.GetAll()
	if err != nil {
		return nil, err
	}
	var musicInfos []MusicInfo
	for _, music := range musics {
		fileurl, _ := s.storage.GetURL(music.FileURL)
		musicinfo := MusicInfo{
			Author:   music.Author,
			Title:    music.Title,
			Duration: music.Duration,
			FileURL:  fileurl,
		}
		musicInfos = append(musicInfos, musicinfo)
	}
	return musicInfos, nil
}

// 获取mp3文件时长
func GetMP3Duration(r io.ReadSeeker) (float64, error) {
	// 保存当前位置
	pos, _ := r.Seek(0, io.SeekCurrent)
	// 无论如何，退出前恢复原来的位置
	defer r.Seek(pos, io.SeekStart)

	// 移动到开头
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	d := mp3.NewDecoder(r)
	var totalDuration time.Duration
	var f mp3.Frame
	var skipped int

	for {
		// 解码下一帧
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		// 累加每一帧的时长
		totalDuration += f.Duration()
	}

	return totalDuration.Seconds(), nil
}

// 上传音乐方法。返回文件完整url
func (s *MusicService) UploadMusic(ctx context.Context, userID uint, file multipart.File,
	fileHeader *multipart.FileHeader, author, title string) (string, error) {
	//验证文件类型
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "audio/") {
		return "", errors.New("只允许上传音频文件")
	}

	if author == "" {
		author = "未知艺术家"
	}

	if title == "" {
		title = "未知标题"
	}

	// 确保从头开始读取
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("预重置文件指针失败: %v\n", err)
	}

	duration, err := GetMP3Duration(file)
	if err != nil {
		duration = 0 // 如果无法获取时长，默认设置为0
		// 可能是其他格式，也不报错，继续上传
	}

	// 再次确保从头开始（给保存文件用）
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("重置文件指针失败: %w", err)
	}

	// 生成文件路径，格式为 audios/{userID}_{timestamp}{ext}
	ext := filepath.Ext(fileHeader.Filename)
	path := fmt.Sprintf("audios/%d_%d%s", userID, time.Now().Unix(), ext)

	// 保存文件
	url, err := s.storage.Upload(ctx, path, file, fileHeader.Size, contentType)
	if err != nil {
		return "", fmt.Errorf("文件存储失败: %w", err)
	}

	err = s.musicRepo.CreateMusic(author, title, int(duration), url, userID)
	if err != nil {
		return "", fmt.Errorf("数据库记录失败: %w", err)
	}

	// 返回音乐URL
	fullURL, err := s.storage.GetURL(url)
	if err != nil {
		return "", fmt.Errorf("获取文件完整URL失败: %w", err)
	}

	return fullURL, nil
}
