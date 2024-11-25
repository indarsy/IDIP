package models

import (
	"time"
)

// 视频模型
type Video struct {
	BaseModel
	FileName   string    `json:"fileName" gorm:"type:varchar(255);not null"`
	FilePath   string    `json:"filePath" gorm:"type:varchar(255);not null"`
	FileSize   int64     `json:"fileSize" gorm:"type:bigint"`
	Duration   float64   `json:"duration" gorm:"type:decimal(10,2)"` // 视频时长(秒)
	WorkshopID uint      `json:"workshopId" gorm:"index"`
	Workshop   Workshop  `json:"workshop" gorm:"foreignKey:WorkshopID"`
	CaptureID  uint      `json:"captureId" gorm:"index"` // 关联的采集任务ID
	StartTime  time.Time `json:"startTime" gorm:"index"`
	EndTime    time.Time `json:"endTime" gorm:"index"`
	Status     int       `json:"status" gorm:"type:tinyint;default:1"` // 1:正常 2:已删除
	Notes      string    `json:"notes" gorm:"type:text"`
}

// 视频查询参数
type VideoQuery struct {
	PageRequest
	WorkshopID uint      `form:"workshopId"`
	StartTime  time.Time `form:"startTime" time_format:"2006-01-02 15:04:05"`
	EndTime    time.Time `form:"endTime" time_format:"2006-01-02 15:04:05"`
	Status     int       `form:"status"`
	Keyword    string    `form:"keyword"`
}

// 视频创建请求
type VideoCreateRequest struct {
	WorkshopID uint      `json:"workshopId" binding:"required"`
	StartTime  time.Time `json:"startTime" binding:"required"`
	EndTime    time.Time `json:"endTime" binding:"required"`
	Notes      string    `json:"notes"`
}

// 视频更新请求
type VideoUpdateRequest struct {
	FileName string `json:"fileName"`
	Notes    string `json:"notes"`
	Status   int    `json:"status"`
}
