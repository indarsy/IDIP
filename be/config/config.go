package config

import (
	"fmt"
	"time"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	RTSP     RTSPConfig     `mapstructure:"rtsp"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
}

type StorageConfig struct {
	Type      string `mapstructure:"type"` // local, s3, oss等
	VideoPath string `mapstructure:"video_path"`
	TempPath  string `mapstructure:"temp_path"`

	// S3配置
	S3 struct {
		Endpoint        string `mapstructure:"endpoint"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		BucketName      string `mapstructure:"bucket_name"`
		Region          string `mapstructure:"region"`
	} `mapstructure:"s3"`
}

type RTSPConfig struct {
	FFmpegPath    string        `mapstructure:"ffmpeg_path"`
	FFprobePath   string        `mapstructure:"ffprobe_path"`
	Timeout       time.Duration `mapstructure:"timeout"`
	SegmentLength int           `mapstructure:"segment_length"` // 视频分段长度(秒)
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`    // 每个日志文件的最大尺寸(MB)
	MaxBackups int    `mapstructure:"max_backups"` // 保留的旧文件最大数量
	MaxAge     int    `mapstructure:"max_age"`     // 保留旧文件的最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩旧文件
}

var GlobalConfig Config

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.Charset,
	)
}
