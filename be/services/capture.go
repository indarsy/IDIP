package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"videodb/be/models"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
)

type CaptureService struct {
	db *gorm.DB
}

func NewCaptureService(db *gorm.DB) *CaptureService {
	return &CaptureService{db: db}
}

func (s *CaptureService) Create(capture *models.Capture) error {
	// 验证时间
	if capture.StartTime.After(capture.EndTime) {
		return fmt.Errorf("开始时间不能晚于结束时间")
	}
	if capture.StartTime.Before(time.Now()) {
		return fmt.Errorf("开始时间不能早于当前时间")
	}

	// 验证间隔
	if capture.Interval < 1 || capture.Interval > 1440 {
		return fmt.Errorf("采集间隔必须在1-1440分钟之间")
	}

	// 设置初始状态
	capture.Status = "waiting"
	if err := s.db.Create(capture).Error; err != nil {
		return err
	}

	// 启动采集任务
	go s.startCapture(capture)

	return nil
}

func (s *CaptureService) List(workshopID uint) ([]models.Capture, error) {
	var captures []models.Capture
	err := s.db.Preload("Workshop").
		Where("workshop_id = ?", workshopID).
		Order("created_at DESC").
		Find(&captures).Error
	return captures, err
}

func (s *CaptureService) ListAll() ([]models.Capture, error) {
	var captures []models.Capture
	err := s.db.Preload("Workshop").Find(&captures).Error
	return captures, err
}

func (s *CaptureService) Get(id uint) (*models.Capture, error) {
	var capture models.Capture
	err := s.db.Preload("Workshop").First(&capture, id).Error
	if err != nil {
		return nil, err
	}
	return &capture, nil
}

func (s *CaptureService) startCapture(capture *models.Capture) {
	// 等待到开始时间
	time.Sleep(time.Until(capture.StartTime))

	// 更新状态为运行中
	s.db.Model(capture).Update("status", "running")

	// 获取车间信息
	var workshop models.Workshop
	if err := s.db.First(&workshop, capture.WorkshopID).Error; err != nil {
		s.updateCaptureStatus(capture, "failed", fmt.Sprintf("获取车间信息失败: %v", err))
		return
	}

	// 创建基础存储目录
	baseStorageDir := filepath.Join("/tmp", "videodb", "storage", "captures")

	// 创建输出目录 - 使用车间ID
	outputDir := filepath.Join(baseStorageDir, fmt.Sprintf("%d", workshop.ID))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		s.updateCaptureStatus(capture, "failed", fmt.Sprintf("创建输出目录失败: %v", err))
		return
	}

	// 计算采集次数
	duration := capture.EndTime.Sub(capture.StartTime)
	intervalDuration := time.Duration(capture.Interval) * time.Minute
	captureCount := int(duration / intervalDuration)

	// 开始采集循环
	currentTime := capture.StartTime
	for i := 0; i < captureCount; i++ {
		// 检查是否已经超过结束时间
		if time.Now().After(capture.EndTime) {
			break
		}

		// 生成输出文件名
		timestamp := currentTime.Format("20060102_150405")
		outputFile := filepath.Join(outputDir, fmt.Sprintf("capture_%s.mp4", timestamp))

		// 执行视频采集
		err := s.captureVideo(workshop.RTSPUrl, outputFile, capture.Interval)
		if err != nil {
			s.updateCaptureStatus(capture, "failed", fmt.Sprintf("视频采集失败: %v", err))
			return
		}

		// 获取文件大小
		fileInfo, err := os.Stat(outputFile)
		if err != nil {
			s.updateCaptureStatus(capture, "failed", fmt.Sprintf("获取文件信息失败: %v", err))
			return
		}

		// 创建视频记录
		video := &models.Video{
			FileName:   filepath.Base(outputFile),
			FilePath:   outputFile,
			FileSize:   fileInfo.Size(),
			Duration:   float64(capture.Interval * 60), // 转换为秒
			WorkshopID: capture.WorkshopID,
			CaptureID:  capture.ID,
			StartTime:  currentTime,
			EndTime:    currentTime.Add(time.Duration(capture.Interval) * time.Minute),
			Status:     1,
			Notes:      fmt.Sprintf("自动采集 - 任务ID:%d", capture.ID),
		}

		if err := s.db.Create(video).Error; err != nil {
			s.updateCaptureStatus(capture, "failed", fmt.Sprintf("保存视频记录失败: %v", err))
			return
		}

		// 更新下一次采集的开始时间
		currentTime = currentTime.Add(intervalDuration)

		// 等待到下一个采集时间点
		if i < captureCount-1 { // 不是最后一次采集
			timeToWait := time.Until(currentTime)
			if timeToWait > 0 {
				time.Sleep(timeToWait)
			}
		}
	}

	// 更新状态为完成
	s.updateCaptureStatus(capture, "completed", "")
}

func (s *CaptureService) captureVideo(rtspUrl string, outputFile string, durationMinutes int) error {
	// 使用 FFmpeg 采集视频
	// 设置采集时长为指定的分钟数
	err := ffmpeg.Input(rtspUrl, ffmpeg.KwArgs{
		"rtsp_transport": "tcp",
		"t":              fmt.Sprintf("%d", durationMinutes*60), // 转换为秒
	}).
		Output(outputFile, ffmpeg.KwArgs{
			"c:v":    "libx264",
			"preset": "medium",
			"crf":    "23",
		}).
		OverWriteOutput().
		Run()

	return err
}

func (s *CaptureService) updateCaptureStatus(capture *models.Capture, status string, message string) {
	updates := map[string]interface{}{
		"status": status,
	}
	if message != "" {
		updates["error_message"] = message
	}
	s.db.Model(capture).Updates(updates)
}

func (s *CaptureService) Cancel(id uint) error {
	return s.db.Model(&models.Capture{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func sanitizeWorkshopName(name string) string {
	// 替换不合法的文件路径字符
	reg := regexp.MustCompile(`[\\/:*?"<>|]`)
	return reg.ReplaceAllString(name, "_")
}
