package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() error {
	configPath := getConfigPath()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	// 读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// 验证配置
	if err := validateConfig(); err != nil {
		return fmt.Errorf("config validation failed: %v", err)
	}

	return nil
}

func getConfigPath() string {
	// 优先使用环境变量中的配置路径
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// 获取当前执行文件的目录
	executable, err := os.Executable()
	if err != nil {
		return "."
	}

	return filepath.Dir(executable)
}

func validateConfig() error {
	// 验证服务器配置
	if GlobalConfig.Server.Port <= 0 {
		return fmt.Errorf("invalid server port")
	}

	// 验证数据库配置
	if GlobalConfig.Database.Port <= 0 {
		return fmt.Errorf("invalid database port")
	}

	// 验证存储路径
	if GlobalConfig.Storage.Type == "local" {
		if err := validatePath(GlobalConfig.Storage.VideoPath); err != nil {
			return fmt.Errorf("invalid video path: %v", err)
		}
		if err := validatePath(GlobalConfig.Storage.TempPath); err != nil {
			return fmt.Errorf("invalid temp path: %v", err)
		}
	}

	return nil
}

func validatePath(path string) error {
	// 检查路径是否存在，不存在则创建
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
