package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 检查并创建目录
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// 生成唯一文件名
func GenerateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 使用时间戳和MD5生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	hash := md5.New()
	hash.Write([]byte(nameWithoutExt + timestamp))
	hashString := hex.EncodeToString(hash.Sum(nil))[:8]

	return fmt.Sprintf("%s_%s%s", timestamp, hashString, ext)
}

// 保存上传文件
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// 获取文件大小
func GetFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// 检查文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 删除文件
func DeleteFile(filename string) error {
	if FileExists(filename) {
		return os.Remove(filename)
	}
	return nil
}
