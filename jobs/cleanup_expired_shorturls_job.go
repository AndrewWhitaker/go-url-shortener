package jobs

import (
	"log"
	"time"
	"url-shortener/models"
	"url-shortener/services"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func CleanupExpiredShortUrls(db *gorm.DB, clock services.Clock) (int64, error) {
	now := clock.Now()
	deleteResult := db.Where("expires_on <= ?", now).Delete(&models.ShortUrl{})

	if deleteResult.Error != nil {
		return 0, deleteResult.Error
	}

	return deleteResult.RowsAffected, nil
}

func StartScheduler(gormDB *gorm.DB, clock services.Clock) {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Seconds().Do(func() {
		deletions, err := CleanupExpiredShortUrls(gormDB, services.SystemClock{})

		if err != nil {
			log.Printf("encountered error deleting expired short urls: %v", err)
			return
		}

		if deletions > 0 {
			log.Printf("deleted %d expired short urls", deletions)
		}
	})

	scheduler.StartAsync()
}
