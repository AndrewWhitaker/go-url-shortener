package services

import (
	"errors"
	"net/url"
	"url-shortener/enums"
	"url-shortener/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/exp/slices"
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
	if request.Slug == "" {
		randomSlug, err := GenerateSlug()

		if err != nil {
			return CreationResult{
				Error: err,
			}
		}
		request.Slug = randomSlug
	}

	validUrl, err := validateLongUrl(request.LongUrl)

	if !validUrl {
		return CreationResult{
			Status: enums.CreationResultInvalidLongUrl,
			Error:  err,
		}
	}

	err = s.DB.Create(&request).Error

	if err == nil {
		return CreationResult{
			Status: enums.CreationResultCreated,
			Record: request,
		}
	}

	if !isUniqueConstraintViolation(err) {
		return CreationResult{
			Error: err,
		}
	}

	var existing models.ShortUrl

	whereClause := models.ShortUrl{}
	whereClause.LongUrl = request.LongUrl

	err = s.DB.
		Where(&whereClause).
		First(&existing).Error

	if err == nil {
		return CreationResult{
			Status: enums.CreationResultAlreadyExists,
			Record: &existing,
		}
	}

	if err == gorm.ErrRecordNotFound {
		return CreationResult{
			Status: enums.CreationResultDuplicateSlug,
		}
	}

	return CreationResult{
		Status: enums.CreationResultUnknownError,
	}
}

func isUniqueConstraintViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}

func validateLongUrl(longUrl string) (bool, error) {
	u, err := url.Parse(longUrl)

	if err != nil {
		return false, err
	}

	validScheme := slices.Contains([]string{"http", "https"}, u.Scheme)

	return validScheme, nil
}
