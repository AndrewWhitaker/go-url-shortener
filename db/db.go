// https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/
package db

import (
	"database/sql"
	"url-shortener/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(sqlDB *sql.DB) (*gorm.DB, error) {
	db, err := ConnectDatabaseWithoutMigrating(sqlDB)

	if err != nil {
		return db, err
	}

	db.AutoMigrate(&models.ShortUrl{}, models.Click{})

	return db, err
}

func ConnectDatabaseWithoutMigrating(sqlDB *sql.DB) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	return db, err
}

type UniqueConstraintName int

const (
	None UniqueConstraintName = iota
	DuplicateLongUrl
	DuplicateSlug
)

func (u UniqueConstraintName) String() string {
	return []string{"uq_short_urls_long_url", "uq_short_urls_slug"}[u]
}

func ParseString(s string) UniqueConstraintName {
	constraintsMap := map[string]UniqueConstraintName{
		"uq_short_urls_long_url": DuplicateLongUrl,
		"uq_short_urls_slug":     DuplicateSlug,
	}

	u, ok := constraintsMap[s]

	if !ok {
		return None
	}

	return u
}
