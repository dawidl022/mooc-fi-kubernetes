package server

import (
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/config"
	"github.com/dawidl022/mooc-fi-kubernetes/todo/models"
)

func initDB(conf *config.Conf) (*gorm.DB, error) {
	db, err := connectDB(conf)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB(conf *config.Conf) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: url.QueryEscape(conf.DatabaseUrl),
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}
