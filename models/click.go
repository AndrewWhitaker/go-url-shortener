package models

import "time"

type Click struct {
	Id         int64 `gorm:"primaryKey"`
	ShortUrlId int64
	CreatedAt  time.Time `gorm:"index:idx_clicks_created_at,sort:asc"`
}
