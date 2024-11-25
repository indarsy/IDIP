package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

// 验证邮箱格式
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 验证手机号格式
func ValidatePhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// 验证URL格式
func ValidateURL(url string) bool {
	pattern := `^(http|https|rtsp)://[^\s/$.?#].[^\s]*$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(url)
}

// 验证文件扩展名
func ValidateFileExt(filename string, allowedExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range allowedExts {
		if ext == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}

// 验证视频文件扩展名
func ValidateVideoExt(filename string) bool {
	allowedExts := []string{".mp4", ".avi", ".mkv", ".mov", ".wmv"}
	return ValidateFileExt(filename, allowedExts)
}
