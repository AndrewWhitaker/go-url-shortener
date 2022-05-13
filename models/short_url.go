package models

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type ShortUrl struct {
	Id     int64   `json:"-"          gorm:"primaryKey"`
	Clicks []Click `json:"-"          gorm:"constraint:OnDelete:CASCADE"`
	ShortUrlReadFields
}

type ShortUrlCreateFields struct {
	LongUrl   string    `json:"long_url"   gorm:"index:uq_short_urls_long_url,unique;not null" binding:"required,url" example:"http://www.google.com" format:"url"`
	ExpiresOn null.Time `json:"expires_on" format:"dateTime" example:"2023-01-01T16:30:00Z"`
	Slug      string    `json:"slug"       gorm:"index:uq_short_urls_slug,unique;not null"  example:"myslug" binding:""`
}

type ShortUrlReadFields struct {
	ShortUrlCreateFields
	CreatedAt time.Time `json:"created_at" format:"dateTime" example:"2022-05-11T11:30:00Z"`
}
