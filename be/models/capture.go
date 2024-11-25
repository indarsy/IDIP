package models

import (
	"time"
)

type Capture struct {
	BaseModel
	WorkshopID   uint      `json:"workshopId" gorm:"not null"`
	StartTime    time.Time `json:"startTime" gorm:"not null"`
	EndTime      time.Time `json:"endTime" gorm:"not null"`
	Interval     int       `json:"interval" gorm:"not null"` // 采集间隔(分钟)
	Status       string    `json:"status"`                   // waiting, running, completed, failed, cancelled
	ErrorMessage string    `json:"errorMessage"`             // 错误信息
	Workshop     Workshop  `json:"workshop" gorm:"foreignKey:WorkshopID"`
}
