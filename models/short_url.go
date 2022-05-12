package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ShortUrl struct {
	Id        int64     `json:"-"          gorm:"primaryKey"`
	LongUrl   string    `json:"long_url"   gorm:"index:uq_short_urls_long_url,unique;not null" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresOn null.Time `json:"expires_on"`
	Slug      string    `json:"slug"       gorm:"index:uq_short_urls_slug,unique;not null"`
	Clicks    []Click   `json:"-"          gorm:"constraint:OnDelete:CASCADE"`
}
