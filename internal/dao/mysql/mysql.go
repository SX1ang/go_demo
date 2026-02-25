package mysql

import (
	"demoProject/config"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host, cfg.Mysql.Port, cfg.Mysql.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("failed to connect database", zap.Error(err))
		return nil, err
	}

	return db, nil
}
