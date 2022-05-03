package models

import "time"

type ShortUrl struct {
	Id        int64  `gorm:"primaryKey"`
	LongUrl   string `gorm:"index:uq_short_urls_long_url,unique" json:"long_url"`
	CreatedAt time.Time
}
