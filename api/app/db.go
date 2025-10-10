package app

import (
	"fybe/model/domain"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(logger *zap.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})

	if err != nil {
		logger.Fatal("failed to connect DB", zap.Error(err))
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Ticket{}, &domain.Registrant{}); err != nil {
		logger.Fatal("auto migrate failed", zap.Error(err))
	}

	logger.Info("DB connected successfully")

	return db
}
