package server

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/dawidl022/mooc-fi-kubernetes/ping-pong/config"
	"github.com/dawidl022/mooc-fi-kubernetes/ping-pong/models"
)

func initDB(conf *config.Conf) (*gorm.DB, error) {
	db, err := connectDB(conf)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Counter{})
	if err != nil {
		return nil, err
	}

	var rowsPresent int64
	err = db.Model(&models.Counter{}).Count(&rowsPresent).Error
	if err != nil {
		return nil, err
	}

	if rowsPresent < 1 {
		err := db.Create(&models.Counter{}).Error
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func connectDB(conf *config.Conf) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: conf.DatabaseUrl,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func getCountFromDB(db *gorm.DB) int {
	var counter models.Counter
	err := db.Find(&counter).Error

	if err != nil {
		return 0
	}
	return counter.PingCount
}

func incrementRequestCountInDB(db *gorm.DB) {
	if db == nil {
		return
	}
	db.Transaction(func(tx *gorm.DB) error {
		return db.Exec("UPDATE counters SET ping_count = ping_count + 1").Error
	})
}
