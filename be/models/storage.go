package models

// 存储配置模型
type Storage struct {
	BaseModel
	Type      string `json:"type" gorm:"type:varchar(20);not null"` // local, s3, oss
	Name      string `json:"name" gorm:"type:varchar(100);not null"`
	Config    string `json:"config" gorm:"type:text"` // JSON格式的配置信息
	IsDefault bool   `json:"isDefault" gorm:"default:false"`
	Status    int    `json:"status" gorm:"type:tinyint;default:1"` // 1:启用 2:禁用
}

// S3存储配置
type S3Config struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	BucketName      string `json:"bucketName"`
	Region          string `json:"region"`
}

// OSS存储配置
type OSSConfig struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	BucketName      string `json:"bucketName"`
}

// 本地存储配置
type LocalConfig struct {
	BasePath string `json:"basePath"`
}
