package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Resources struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewResources(db *gorm.DB, redis *redis.Client) (*Resources, func(), error) {
	res := &Resources{
		db:    db,
		redis: redis,
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
