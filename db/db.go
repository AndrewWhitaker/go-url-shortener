// https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
package db

import (
	"database/sql"
	"url-shortener/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(sqlDB *sql.DB) error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	db.AutoMigrate(&models.ShortUrl{})

	DB = db

	return err
}
