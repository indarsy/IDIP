package services

import (
	"errors"
	"fmt"
	"videodb/be/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkshopService struct {
	db *gorm.DB
}

func NewWorkshopService(db *gorm.DB) *WorkshopService {
	return &WorkshopService{db: db}
}

// 获取车间列表
func (s *WorkshopService) List() ([]models.Workshop, error) {
	var workshops []models.Workshop
	err := s.db.Find(&workshops).Error
	return workshops, err
}

// 创建车间
func (s *WorkshopService) Create(workshop *models.Workshop) error {
	return s.db.Clauses(clause.OnConflict{
		UpdateAll: true, // 冲突时更新所有字段
	}).Create(workshop).Error
}

// 更新车间信息
func (s *WorkshopService) Update(id uint, workshop *models.Workshop) error {
	return s.db.Model(&models.Workshop{}).Where("id = ?", id).Updates(workshop).Error
}

// 删除车间
func (s *WorkshopService) Delete(id uint) error {
	return s.db.Delete(&models.Workshop{}, id).Error
}

// 更新车间状态
func (s *WorkshopService) UpdateStatus(id uint, status int) error {
	return s.db.Model(&models.Workshop{}).Where("id = ?", id).Update("status", status).Error
}

// GetByID 根据ID获取车间信息
func (s *WorkshopService) GetByID(id uint) (*models.Workshop, error) {
	var workshop models.Workshop

	// 使用GORM查询数据库
	err := s.db.First(&workshop, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("workshop not found with id: %d", id)
		}
		return nil, fmt.Errorf("failed to get workshop: %v", err)
	}

	return &workshop, nil
}
