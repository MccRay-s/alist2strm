package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/MccRay-s/alist2strm/config"
	"github.com/MccRay-s/alist2strm/model/configs"
	"github.com/MccRay-s/alist2strm/model/filehistory"
	"github.com/MccRay-s/alist2strm/model/invalidstrm"
	"github.com/MccRay-s/alist2strm/model/notification"
	"github.com/MccRay-s/alist2strm/model/task"
	"github.com/MccRay-s/alist2strm/model/tasklog"
	"github.com/MccRay-s/alist2strm/model/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.AppConfig) error {
	dbDir := cfg.Database.BaseDir
	dbPath := filepath.Join(dbDir, cfg.Database.Name)

	// 创建数据库目录
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	DB = db
	log.Printf("数据库连接成功: %s", dbPath)

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(
		&user.User{},
		&configs.Config{},
		&task.Task{},
		&tasklog.TaskLog{},
		&filehistory.FileHistory{},
		&notification.Queue{},
		&invalidstrm.InvalidStrmFile{},
	); err != nil {
		return fmt.Errorf("数据库表迁移失败: %v", err)
	}
	log.Printf("数据库表结构同步完成")

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
