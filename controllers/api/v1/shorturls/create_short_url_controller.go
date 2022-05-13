package shorturls

import (
	"net/http"
	"url-shortener/e"
	"url-shortener/enums"
	"url-shortener/middleware"
	"url-shortener/models"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
)

type CreateShortUrlController struct {
	CreateShortUrlService *services.CreateShortUrlService
}

// CreateShortUrl godoc
// @Summary      Create a new short url
// @Description  Create a new short url. Users may specify a slug and an expiration date. If a slug is not supplied, an 8 character slug will automatically be generated for the short url.
// @Tags         shorturls
// @Accept       json
// @Produce      json
// @Param        shorturl  body      models.ShortUrlCreateFields  true  "New short URL parameters"
// @Success      200       {object}  models.ShortUrlReadFields
// @Success      201       {object}  models.ShortUrlReadFields
// @Failure      400       {object}  e.ErrorResponse
// @Failure      404
// @Failure      409  {object}  e.ErrorResponse
// @Failure      500
// @Router       /shorturls [post]
func (controller *CreateShortUrlController) HandleRequest(c *gin.Context, request models.ShortUrl) {
	createResult := controller.CreateShortUrlService.Create(&request)

	if createResult.Error != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var status int
	var body interface{}

	switch createResult.Status {
	case enums.CreationResultCreated:
		status = http.StatusCreated
		body = createShortUrlResponseHelper{
			Host:     c.Request.Host,
			ShortUrl: *createResult.Record,
		}
	case enums.CreationResultAlreadyExists:
		status = http.StatusOK
		body = createShortUrlResponseHelper{
			Host:     c.Request.Host,
			ShortUrl: *createResult.Record,
		}
	case enums.CreationResultDuplicateSlug:
		status = http.StatusConflict
		body = e.ErrorResponse{
			Errors: []e.ValidationError{
				e.ValidationError{
					Field:  "Slug",
					Reason: "must be unique",
				},
			},
		}
	}

	c.JSON(status, body)
}

func (controller *CreateShortUrlController) Register(r *gin.Engine) {
	r.POST("/api/v1/shorturls", middleware.ModelBindingWrapper[models.ShortUrl](controller))
}

type CreateShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
	models.ShortUrlReadFields
}
