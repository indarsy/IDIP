package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type RTSPService struct {
	recordings     map[uint]*exec.Cmd
	recordingMutex sync.Mutex
}

func NewRTSPService() *RTSPService {
	return &RTSPService{
		recordings: make(map[uint]*exec.Cmd),
	}
}

// 开始录制
func (s *RTSPService) StartRecording(ctx context.Context, rtspURL string, outputPath string, workshopID uint) error {
	s.recordingMutex.Lock()
	defer s.recordingMutex.Unlock()

	// 检查是否已经在录制
	if _, exists := s.recordings[workshopID]; exists {
		return fmt.Errorf("workshop %d is already recording", workshopID)
	}

	// 使用FFmpeg录制RTSP流
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", rtspURL,
		"-c", "copy",
		"-f", "mp4",
		outputPath)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start recording: %v", err)
	}

	s.recordings[workshopID] = cmd

	// 启动goroutine监控录制进程
	go func() {
		defer func() {
			s.recordingMutex.Lock()
			delete(s.recordings, workshopID)
			s.recordingMutex.Unlock()
		}()

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Recording for workshop %d ended with error: %v\n", workshopID, err)
		}
	}()

	return nil
}

// 停止录制
func (s *RTSPService) StopRecording(workshopID uint) error {
	s.recordingMutex.Lock()
	defer s.recordingMutex.Unlock()

	cmd, exists := s.recordings[workshopID]
	if !exists {
		return fmt.Errorf("no recording found for workshop %d", workshopID)
	}

	// 发送中断信号给FFmpeg进程
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to stop recording: %v", err)
	}

	return nil
}

// 获取预览流地址
func (s *RTSPService) GetPreviewURL(workshopID uint) (string, error) {
	// 这里可以实现转换RTSP流为HLS或WebRTC流的逻辑
	// 示例中仅返回模拟地址
	return fmt.Sprintf("http://localhost:8882/live/%d/index.m3u8", workshopID), nil
}

// 检查RTSP流是否可用
func (s *RTSPService) CheckRTSPStream(rtspURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffprobe",
		"-i", rtspURL,
		"-v", "quiet")

	return cmd.Run()
}

// 获取录制状态
func (s *RTSPService) GetRecordingStatus(workshopID uint) bool {
	s.recordingMutex.Lock()
	defer s.recordingMutex.Unlock()

	_, exists := s.recordings[workshopID]
	return exists
}
