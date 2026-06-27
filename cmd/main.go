package main

import (
	"blog/config"
	"blog/routes"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	// 加载环境变量
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		logrus.Info("Failed to resolve current file path")
	} else {
		rootDir := filepath.Dir(filepath.Dir(currentFile))
		envPath := filepath.Join(rootDir, ".env")
		if err := godotenv.Load(envPath); err != nil {
			logrus.Warn("No .env file found, using system enviroment variables")
		}
	}

	// 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting blog application...")

	// 创建文件上传目录
	var uploadPath string
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	if config.GetEnv("GIN_MODE", "release") == "debug" {
		uploadPath = filepath.Join(rootDir, "uploads")
	} else {
		uploadPath = config.GetEnv("UPLOAD_DIR", filepath.Join(rootDir, "uploads"))
		if !filepath.IsAbs(uploadPath) {
			uploadPath = filepath.Join(rootDir, uploadPath)
		}
	}

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logrus.Warn("Failed to create uploads directory")
	}

	// 初始化数据库
	//config.InitDataBase()
	// 使用MySQL
	config.InitMysqlDataBase()

	// 路由设置
	r := routes.SetupRoutes(uploadPath)

	// 启动服务器
	port := config.GetEnv("PORT", "8080")
	port = ":" + port
	logrus.WithField("port", port).Info("Server Starting")

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}
}
