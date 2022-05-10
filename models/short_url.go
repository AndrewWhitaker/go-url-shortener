package models

import "time"

type ShortUrl struct {
	Id        int64  `gorm:"primaryKey"`
	LongUrl   string `gorm:"index:uq_short_urls_long_url,unique;not null" json:"long_url" binding:"required"`
	CreatedAt time.Time
	Slug      string  `gorm:"index:uq_short_urls_slug,unique;not null", json:"slug" binding:`
	Clicks    []Click `gorm:"constraint:OnDelete:CASCADE"`
}
