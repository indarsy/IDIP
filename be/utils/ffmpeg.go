package utils

import (
	"context"
	"os/exec"
	"time"
)

// FFmpeg工具结构体
type FFmpeg struct {
	BinPath string
}

// 创建新的FFmpeg实例
func NewFFmpeg(binPath string) *FFmpeg {
	if binPath == "" {
		binPath = "ffmpeg"
	}
	return &FFmpeg{BinPath: binPath}
}

// 检查FFmpeg是否可用
func (f *FFmpeg) CheckAvailable() error {
	cmd := exec.Command(f.BinPath, "-version")
	return cmd.Run()
}

// 获取视频信息
func (f *FFmpeg) GetVideoInfo(filepath string) (map[string]string, error) {
	cmd := exec.Command(f.BinPath,
		"-i", filepath,
		"-show_format",
		"-show_streams",
		"-v", "quiet",
		"-print_format", "json",
	)

	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析输出信息
	info := make(map[string]string)
	// 这里需要根据实际输出格式解析JSON
	return info, nil
}

// 转换视频格式
func (f *FFmpeg) ConvertVideo(ctx context.Context, input, output string, options map[string]string) error {
	args := []string{"-i", input}

	// 添加转换选项
	for k, v := range options {
		args = append(args, "-"+k, v)
	}

	args = append(args, output)
	cmd := exec.CommandContext(ctx, f.BinPath, args...)

	return cmd.Run()
}

// 生成视频缩略图
func (f *FFmpeg) GenerateThumbnail(input, output string, timestamp string) error {
	cmd := exec.Command(f.BinPath,
		"-i", input,
		"-ss", timestamp,
		"-vframes", "1",
		"-vf", "scale=320:-1",
		output,
	)

	return cmd.Run()
}

// 截取视频片段
func (f *FFmpeg) CutVideo(input, output string, start, duration string) error {
	cmd := exec.Command(f.BinPath,
		"-i", input,
		"-ss", start,
		"-t", duration,
		"-c", "copy",
		output,
	)

	return cmd.Run()
}

// 检查RTSP流是否可用
func (f *FFmpeg) CheckRTSPStream(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, f.BinPath,
		"-i", url,
		"-t", "1",
		"-f", "null",
		"-",
	)

	return cmd.Run()
}
