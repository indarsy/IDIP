package config

const (
	// 视频状态
	VideoStatusNormal  = 1
	VideoStatusDeleted = 2

	// 车间状态
	WorkshopStatusOffline = 0
	WorkshopStatusOnline  = 1

	// 录制状态
	RecordingStatusStopped = 0
	RecordingStatusRunning = 1

	// 存储类型
	StorageTypeLocal = "local"
	StorageTypeS3    = "s3"
	StorageTypeOSS   = "oss"

	// 文件类型
	FileTypeVideo = "video"
	FileTypeImage = "image"

	// 默认值
	DefaultPageSize   = 10
	MaxPageSize       = 100
	DefaultTimeFormat = "2006-01-02 15:04:05"

	// 错误码
	ErrCodeSuccess       = 0
	ErrCodeParamInvalid  = 1001
	ErrCodeDBError       = 1002
	ErrCodeFileNotFound  = 1003
	ErrCodeUnauthorized  = 1004
	ErrCodeInternalError = 1005
)
