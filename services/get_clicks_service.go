package services

import (
	"url-shortener/enums"

	"gorm.io/gorm"
)

type GetClicksService struct {
	DB *gorm.DB
}

type GetClicksResult struct {
	Status enums.GetClicksStatus
	Count  int64
	Error  error
}

func (s *GetClicksService) GetClicks(slug string) GetClicksResult {
	var count int64
	countResult := s.DB.Raw(`
			SELECT COUNT(*)
			FROM
				short_urls
				LEFT OUTER JOIN clicks ON clicks.short_url_id = short_urls.id
			WHERE
				short_urls.slug = ?
			GROUP BY short_urls.id
		`, slug).Scan(&count)

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
