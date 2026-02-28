package infra

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Resources struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewResources(db *gorm.DB, redis *redis.Client) (*Resources, func(), error) {
	res := &Resources{
		DB:    db,
		Redis: redis,
	}

	close := func() {
		if db != nil {
			sqlDB, err := db.DB()
			if err == nil {
				sqlDB.Close()
			}
		}
		if redis != nil {
			redis.Close()
		}
	}

	return res, close, nil
}
