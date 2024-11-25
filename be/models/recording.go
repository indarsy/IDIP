package models

import (
	"time"
)

// 录制任务模型
type Recording struct {
	BaseModel
	WorkshopID uint      `json:"workshopId" gorm:"index"`
	Workshop   Workshop  `json:"workshop" gorm:"foreignKey:WorkshopID"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Status     int       `json:"status" gorm:"type:tinyint;default:0"` // 0:未开始 1:进行中 2:已完成 3:已失败
	VideoID    uint      `json:"videoId" gorm:"index"`
	Video      Video     `json:"video" gorm:"foreignKey:VideoID"`
	ErrorMsg   string    `json:"errorMsg" gorm:"type:text"`
}

// 录制任务创建请求
type RecordingCreateRequest struct {
	WorkshopID uint      `json:"workshopId" binding:"required"`
	StartTime  time.Time `json:"startTime" binding:"required"`
	EndTime    time.Time `json:"endTime" binding:"required"`
}

// 录制状态更新请求
type RecordingStatusUpdate struct {
	Status   int    `json:"status" binding:"required,oneof=0 1 2 3"`
	VideoID  uint   `json:"videoId"`
	ErrorMsg string `json:"errorMsg"`
}
