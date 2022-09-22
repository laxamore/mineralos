package db

import (
	"fmt"
	"github.com/laxamore/mineralos/config"
	"github.com/laxamore/mineralos/internal/db/models"
	"github.com/laxamore/mineralos/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type IDB interface {
	Create(value interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
}

func ConnectDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		config.Config.DB_HOST, config.Config.DB_USER, config.Config.DB_PASSWORD, config.Config.DB_NAME, config.Config.DB_PORT)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	DB.DB()

	if err != nil {
		logger.Panic("Failed connecting to db", err)
	}
}

func InitDB() {
	for _, model := range models.GetModels() {
		if err := DB.AutoMigrate(model); err != nil {
			logger.Panic("Failed to migrate model", err)
		}

		switch model.(type) {
		case *models.Role:
			DB.Create(&models.RoleAdmin)
			DB.Create(&models.RoleOperator)
			DB.Create(&models.RoleUser)
		}
	}
}
