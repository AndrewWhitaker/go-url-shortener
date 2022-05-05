package services

import (
	"errors"
	"url-shortener/enums"
	"url-shortener/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

type CreateShortUrlService struct {
	DB *gorm.DB
}

type CreationResult struct {
	Status enums.CreationStatus
	Record *models.ShortUrl
	Error  error
}

func (s *CreateShortUrlService) Create(request *models.ShortUrl) CreationResult {
	err := s.DB.Create(&request).Error

	if err == nil {
		return CreationResult{
			Status: enums.CreationResultCreated,
			Record: request,
		}
	}

	if isUniqueConstraintViolation(err) {
		var existing models.ShortUrl

		s.DB.Where(&models.ShortUrl{LongUrl: request.LongUrl}).First(&existing)

		return CreationResult{
			Status: enums.CreationResultAlreadyExisted,
			Record: &existing,
		}
	}

	return CreationResult{
		Error: err,
	}
}

func isUniqueConstraintViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
