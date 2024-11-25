package main

import (
	"fmt"
	"log"
	"time"
	"videodb/be/config"
	"videodb/be/handlers"
	"videodb/be/models"
	"videodb/be/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化配置
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&config.GlobalConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}
}

// 初始化数据库连接
func initDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GlobalConfig.Database.Username,
		config.GlobalConfig.Database.Password,
		config.GlobalConfig.Database.Host,
		config.GlobalConfig.Database.Port,
		config.GlobalConfig.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	err = db.AutoMigrate(&models.Video{}, &models.Workshop{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// 设置路由
func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源，修改操作需要处理，查询不用
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 使用全局配置
	cfg := &config.GlobalConfig

	// 创建服务实例
	videoService := services.NewVideoService(db)
	rtspService := services.NewRTSPService()
	workshopService := services.NewWorkshopService(db)
	captureService := services.NewCaptureService(db)
	webrtcService := services.NewWebRTCService(cfg)

	// 创建处理器实例
	videoHandler := handlers.NewVideoHandler(videoService, rtspService, workshopService)
	workshopHandler := handlers.NewWorkshopHandler(workshopService, rtspService)
	captureHandler := handlers.NewCaptureHandler(captureService)
	webrtcHandler := handlers.NewWebRTCHandler(webrtcService) // 添加 WebRTC 处理器

	// API 路由组
	api := r.Group("/api") // 设置api前缀
	{
		// 视频相关路由
		videos := api.Group("/videos")
		{
			videos.GET("", videoHandler.List)
			//videos.POST("", videoHandler.CreateVideo)
			videos.GET("/:id", videoHandler.Get)
			//videos.PUT("/:id", videoHandler.UpdateVideo)
			videos.DELETE("/:id", videoHandler.Delete)
			videos.GET("/:id/download", videoHandler.Download)
			videos.DELETE("/batch", videoHandler.BatchDelete)
		}

		// 车间相关路由
		workshops := api.Group("/workshops")
		{
			workshops.GET("", workshopHandler.List)
			workshops.POST("", workshopHandler.Create)
			workshops.PUT("/:id", workshopHandler.Update)
			workshops.DELETE("/:id", workshopHandler.Delete)
		}

		// RTSP 流相关路由
		rtsp := api.Group("/rtsp")
		{
			rtsp.POST("/start", videoHandler.StartRecording)
			rtsp.POST("/stop", videoHandler.StopRecording)
			//rtsp.GET("/preview/:workshopId", videoHandler.PreviewStream)
		}

		// 采集相关路由
		captures := api.Group("/captures")
		{
			captures.POST("", captureHandler.Create)
			captures.GET("", captureHandler.List)
			captures.GET("/:id", captureHandler.Get)
			captures.POST("/:id/cancel", captureHandler.Cancel)
		}

		// WebRTC 相关路由
		webrtc := api.Group("/webrtc")
		{
			webrtc.POST("", webrtcHandler.HandleWebRTC)
		}

	}

	return r
}

// 启动定时任务
func startCronJobs(videoService *services.VideoService) {
	// 每天凌晨清理过期视频
	c := cron.New(cron.WithLocation(time.Local))
	c.AddFunc("0 0 * * *", func() {
		if err := videoService.CleanExpiredVideos(); err != nil {
			log.Printf("Failed to clean expired videos: %v", err)
		}
	})
	c.Start()
}

func main() {
	// 初始化配置
	initConfig()

	// 初始化数据库连接
	db := initDB()

	// 设置路由
	r := setupRouter(db)

	// 启动定时任务
	// videoService := services.NewVideoService(db)
	// startCronJobs(videoService)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	log.Printf("Server starting on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
