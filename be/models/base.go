package models

import (
	"time"
)

// 基础模型
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	//DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // gorm.DeletedAt 类型用于软删除
}

// 分页请求
type PageRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=100"`
}

// 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
