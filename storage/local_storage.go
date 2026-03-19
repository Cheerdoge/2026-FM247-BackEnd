package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage 创建一个新的 LocalStorage 实例
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		fmt.Printf("无法创建本地目录: %v\n", err)
		return nil
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// 上传文件到本地存储, 返回相对路径
func (s *LocalStorage) Upload(ctx context.Context, path string, content io.Reader, fileSize int64, contentType string) (string, error) {
	// 1. 拼接绝对路径 (basePath + path)，例如 "./uploads" + "music/1.mp3"
	fullPath := filepath.Join(s.basePath, path)

	// 2. 确保目标文件夹存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 3. 创建文件
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 4. 写入内容
	if _, err = io.Copy(out, content); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 5. 返回相对路径
	return path, nil
}

// 删除文件
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// 相对路径转换为完整路径
func (s *LocalStorage) GetURL(path string) (string, error) {
	// 修正路径分隔符问题，确保 URL 是正斜杠
	urlPath := filepath.ToSlash(path)
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	return s.baseURL + urlPath, nil
}
