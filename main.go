package main

import (
	"2026-FM247-BackEnd/config"
	handler "2026-FM247-BackEnd/handlers"
	"2026-FM247-BackEnd/logger"
	repository "2026-FM247-BackEnd/repositories"
	"2026-FM247-BackEnd/router"
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/storage"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 优先从 go.env 加载环境变量（开发环境）
	_ = godotenv.Overload("go.env")

	fmt.Println("正在初始化配置")
	config.LoadConfig()

	logger.InitLogger("INFO")

	db, err := config.ConnectDatabase()
	if err != nil {
		fmt.Printf("无法连接到mysql数据库: %v\n", err)
		return
	}

	redisClient, err := config.ConnectRedis()
	if err != nil {
		fmt.Printf("无法连接到Redis: %v\n", err)
		return
	}

	storage := storage.NewLocalStorage("./uploads", config.AppConfig.BaseURL)

	fmt.Println("数据库连接成功")

	// aiclient初始化
	aiClient, err := config.NewAIClient()
	if err != nil {
		fmt.Printf("无法创建AI客户端: %v\n", err)
		return
	}

	//dao层初始化
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenBlacklistRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	studyDataRepo := repository.NewStudyDataRepository(db, redisClient)
	musicRepo := repository.NewMusicRepository(db)
	ambientSoundRepo := repository.NewAmbientSoundRepository(db)
	aichatRepo := repository.NewAIChatRepository(redisClient)

	//service层初始化
	userService := service.NewUserService(userRepo, tokenRepo, storage)
	tokenService := service.NewTokenBlacklistService(tokenRepo)
	todoService := service.NewTodoService(todoRepo)
	musicService := service.NewMusicService(musicRepo, storage)
	studyDataService := service.NewStudyDataService(studyDataRepo)
	ambientSoundService := service.NewAmbientSoundService(ambientSoundRepo, storage)
	aichatService := service.NewAIChatService(aichatRepo, aiClient)
	//handler层初始化
	authhandler := handler.NewAuthHandler(tokenService, userService, studyDataService)
	avatarHandler := handler.NewAvatarHandler(userService)
	todohandler := handler.NewTodoHandler(todoService)
	studydatahandler := handler.NewStudyDataHandler(studyDataService)
	musichandler := handler.NewMusicHandler(musicService)
	ambientSoundHandler := handler.NewAmbientSoundHandler(ambientSoundService)
	aiChatHandler := handler.NewAIChatHandler(aichatService)

	// 启动服务器
	// r := gin.New()
	// r.Use(gin.Recovery())
	// r.Use(middleware.GinLogger())
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ // 允许的请求源
			"http://localhost:5173", // 前端vite的默认启动地址
			"http://localhost:3000", // 前端自己定义的启动地址
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                             // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"}, // 允许的请求头
		AllowCredentials: true,
	}))

	router.RegisterRoutes(r, authhandler, avatarHandler, todohandler, studydatahandler, musichandler, ambientSoundHandler, aiChatHandler)
	port := ":" + config.AppConfig.ServerPort
	fmt.Printf("服务器正在运行，监听端口 %s\n", port)
	if err := r.Run(port); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}

	fmt.Println("正在与数据库断开连接")
	err = config.CloseDatabase(db)
	if err != nil {
		fmt.Printf("关闭数据库连接时出错: %v\n", err)
		return
	}
	fmt.Println("数据库连接已关闭")
}
