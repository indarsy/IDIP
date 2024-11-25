package models

// 车间模型
type Workshop struct {
	BaseModel
	Name        string  `json:"name" gorm:"type:varchar(100);not null;unique"`
	RTSPUrl     string  `json:"rtspUrl" gorm:"type:varchar(255);not null"`
	Status      int     `json:"status" gorm:"type:tinyint;default:0"` // 2:离线 1:在线
	Description string  `json:"description" gorm:"type:text"`
	Videos      []Video `json:"videos" gorm:"foreignKey:WorkshopID"`
}

// 车间创建请求
type WorkshopCreateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	RTSPUrl     string `json:"rtspUrl" binding:"required,url"`
	Description string `json:"description"`
}

// 车间更新请求
type WorkshopUpdateRequest struct {
	Name        string `json:"name" binding:"max=100"`
	RTSPUrl     string `json:"rtspUrl" binding:"url"`
	Status      int    `json:"status" binding:"oneof=0 1"`
	Description string `json:"description"`
}
