package services

import (
	"time"
	"url-shortener/enums"

	"gorm.io/gorm"
)

type GetClicksService struct {
	DB    *gorm.DB
	Clock Clock
}

type GetClicksResult struct {
	Status enums.GetClicksStatus
	Count  int64
	Error  error
}

func (s *GetClicksService) GetClicks(slug string, timePeriod enums.GetClicksTimePeriod) GetClicksResult {

	var query *gorm.DB

	now := s.Clock.Now()

	// These are flawed calculations in anything but UTC, but close
	// enough for now.
	twentyFourHours := time.Duration(time.Hour * 24)
	oneWeek := time.Duration(twentyFourHours * 7)

	switch timePeriod {
	case enums.GetClicksTimePeriodAllTime:
		query = s.AllTimeClicksQuery(slug)
	case enums.GetClicksTimePeriodPastWeek:
		time := now.Add(-oneWeek)
		query = s.ClicksAfter(slug, time)
	case enums.GetClicksTimePeriod24Hours:
		time := now.Add(-twentyFourHours)
		query = s.ClicksAfter(slug, time)
	}

	var count int64
	countResult := query.Scan(&count)

	if countResult.Error != nil {
		return GetClicksResult{
			Error:  countResult.Error,
			Status: enums.GetClicksResultUnknownError,
		}
	}

	if countResult.RowsAffected == 1 {
		return GetClicksResult{
			Count:  count,
			Status: enums.GetClicksResultSuccessful,
		}
	}

	return GetClicksResult{
		Status: enums.GetClicksResultNotFound,
	}
}

func (s *GetClicksService) AllTimeClicksQuery(slug string) *gorm.DB {
	return s.DB.Raw(`
			SELECT COUNT(*)
			FROM
				short_urls
				LEFT OUTER JOIN clicks ON clicks.short_url_id = short_urls.id
			WHERE
				short_urls.slug = ?
			GROUP BY short_urls.id
	`, slug)
}

func (s *GetClicksService) ClicksAfter(slug string, startTime time.Time) *gorm.DB {
	return s.DB.Raw(`
			SELECT COUNT(*)
			FROM
				short_urls
				LEFT OUTER JOIN clicks ON
					clicks.short_url_id = short_urls.id AND
					clicks.created_at >= ?
			WHERE
				short_urls.slug = ?
			GROUP BY short_urls.id
	`, startTime, slug)
}
