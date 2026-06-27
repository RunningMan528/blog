package config

import (
	"blog/config"
	"blog/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

//  InitDataBase 初始化数据库连接
// 默认使用sqlite进行测试

func InitDataBase() {
	// 获取db文件存储目录
	dbDir, err := getDBDir()
	if err != nil {
		log.Fatalf("init database fail : %v", err)
	}

	// 数据库文件存储的具体位置
	dbPath := filepath.Join(dbDir, "blog.db")

	// 加载环境变量查看是否是debug模式
	debug := GetEnv("GIN_MODE", "release")
	logMode := logger.Silent
	if debug == "debug" {
		logMode = logger.Info
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// 日志配置
		Logger: logger.Default.LogMode(logMode),
		//  命名策略: 自定义GORM的表和列命名方式
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // 表名前缀
			SingularTable: false, // 是否使用单数表名(false : User -> users)
			NoLowerCase:   false, // 禁用自动小写转换
			NameReplacer:  nil,   // 自定义名称替换函数
		},
	})
	if err != nil {
		log.Fatalf("init database fail : %v", err)
	}

	DB = db

	// 自动迁移数据库表结构
	migreateErr := DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if migreateErr != nil {
		log.Fatalf("database migreate err : %v", migreateErr)
	}

	log.Println("sqlite database init successfully!")
}

func InitMysqlDataBase() {
	var err error
	// 从环境变量获取MySQL连接配置
	dbHost := GetEnv("DB_HOST", "localhost")
	dbPort := GetEnv("DB_PORT", "3306")
	dbUser := GetEnv("DB_USER", "root")
	dbPassword := GetEnv("DB_PASSWORD", "")
	dbName := GetEnv("DB_NAME", "blog")
	if config.GetEnv("GIN_MODE", "release") == "debug" {
		dbHost = "localhost"
	}

	// 构建MySQL连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 连接MySQL数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}

	// 自动迁移数据库表结构
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("MySQL database connected and migrated successfully")
}

// 获取数据库存储的db文件目录路径
func getDBDir() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrInvalid
	}

	// 当前文件的目录
	currentDir := filepath.Dir(currentFile)
	// 当前目录的目录
	rootDir := filepath.Dir(currentDir)

	// db文件目录
	dbDir := filepath.Join(rootDir, "db")

	// 确保db目录存在
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	return dbDir, nil

}

// getEnv 获取环境变量,如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
