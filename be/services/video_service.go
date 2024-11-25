package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"videodb/be/config"
	"videodb/be/models"

	"gorm.io/gorm"
)

type VideoService struct {
	db *gorm.DB
}

type VideoQuery struct {
	WorkshopID uint
	StartTime  time.Time
	EndTime    time.Time
	Page       int
	PageSize   int
	Preload    []string
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{db: db}
}

// 获取视频列表
func (s *VideoService) List(query VideoQuery) ([]models.Video, int64, error) {
	db := s.db.Model(&models.Video{})

	// 添加预加载
	for _, preload := range query.Preload {
		db = db.Preload(preload)
	}

	// 添加查询条件
	if query.WorkshopID > 0 {
		db = db.Where("workshop_id = ?", query.WorkshopID)
	}
	if !query.StartTime.IsZero() {
		db = db.Where("start_time >= ?", query.StartTime)
	}
	if !query.EndTime.IsZero() {
		db = db.Where("end_time <= ?", query.EndTime)
	}

	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询数据
	var videos []models.Video
	err := db.Offset((query.Page - 1) * query.PageSize).
		Limit(query.PageSize).
		Order("created_at DESC").
		Find(&videos).Error

	return videos, total, err
}

// 创建视频记录
func (s *VideoService) Create(video *models.Video) error {
	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()
	return s.db.Create(video).Error
}

// 更新视频记录
func (s *VideoService) Update(id uint, video *models.Video) error {
	video.UpdatedAt = time.Now()
	return s.db.Model(&models.Video{}).Where("id = ?", id).Updates(video).Error
}

// 删除视频记录
func (s *VideoService) Delete(id uint) error {
	// 先获取视频信息
	var video models.Video
	if err := s.db.First(&video, id).Error; err != nil {
		return err
	}

	// 删除物理文件
	if err := os.Remove(video.FilePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// 删除数据库记录
	return s.db.Delete(&video).Error
}

// 清理过期视频
func (s *VideoService) CleanExpiredVideos() error {
	// 设置过期时间，例如30天
	expireTime := time.Now().AddDate(0, 0, -30)

	var videos []models.Video
	if err := s.db.Where("create_time < ?", expireTime).Find(&videos).Error; err != nil {
		return err
	}

	for _, video := range videos {
		if err := s.Delete(video.ID); err != nil {
			return err
		}
	}

	return nil
}

// 生成视频存储路径
func (s *VideoService) GenerateVideoPath(workshopID uint) string {
	now := time.Now()
	fileName := fmt.Sprintf("%d_%s.mp4", workshopID, now.Format("20060102150405"))
	return filepath.Join(config.GlobalConfig.Storage.VideoPath, fileName)
}

// GetByID 根据ID获取视频信息
func (s *VideoService) GetByID(id uint) (*models.Video, error) {
	var video models.Video

	// 使用GORM查询数据库，同时预加载关联的Workshop信息
	err := s.db.Preload("Workshop").First(&video, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("video not found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to get video: %v", err)
	}

	return &video, nil
}

// BatchDelete 批量删除视频
func (s *VideoService) BatchDelete(ids []uint) error {
	// 先获取所有要删除的视频信息
	var videos []models.Video
	if err := s.db.Where("id IN ?", ids).Find(&videos).Error; err != nil {
		return err
	}

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除数据库记录
		if err := tx.Delete(&models.Video{}, "id IN ?", ids).Error; err != nil {
			return err
		}

		// 删除物理文件
		for _, video := range videos {
			if video.FilePath != "" {
				if err := os.Remove(video.FilePath); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("删除文件失败: %v", err)
				}
			}
		}

		return nil
	})
}
