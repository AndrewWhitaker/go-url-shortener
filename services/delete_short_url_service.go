package services

import (
	"url-shortener/enums"
	"url-shortener/models"

	"gorm.io/gorm"
)

type DeleteShortUrlService struct {
	DB *gorm.DB
}

type DeleteResult struct {
	Status enums.DeleteStatus
	Record *models.ShortUrl
	Error  error
}

func (s *DeleteShortUrlService) Delete(slug string) DeleteResult {
	var shortUrl models.ShortUrl

	res := s.DB.
		Where(&models.ShortUrl{Slug: slug}).
		Delete(&shortUrl)

	response := DeleteResult{}

	if res.Error == nil {
		if res.RowsAffected == 1 {
			response.Status = enums.DeleteResultSuccessful
		} else {
			response.Status = enums.DeleteResultNotFound
		}
	} else {
		response.Status = enums.DeleteResultUnknownError
		response.Error = res.Error
	}

	return response
}
