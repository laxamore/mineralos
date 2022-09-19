package databases

import (
	"fmt"
	"github.com/laxamore/mineralos/internal/databases/models"
	"github.com/laxamore/mineralos/internal/databases/models/user"
	"github.com/laxamore/mineralos/internal/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
)

var (
	DB *gorm.DB
)

type DBInterface interface {
	Create(value interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
}

func Connect(dbUser string, dbPassword string, dbHost string, dbPort int, dbName string) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	DB.DB()

	if err != nil {
		logger.Panic("Failed connecting to databases", err)
	}
}

func InitDatabase() {
	for _, model := range models.GetModels() {
		if err := DB.AutoMigrate(model); err != nil {
			logger.Panicf("Failed to migrate model", err)
		}

		if reflect.TypeOf(model) == reflect.TypeOf(user.Role{}) {
			DB.Create(&[]user.Role{user.RoleAdmin, user.RoleOperator, user.RoleUser})
		}
	}
}
