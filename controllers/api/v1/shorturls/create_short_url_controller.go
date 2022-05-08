package shorturls

import (
	"encoding/json"
	"net/http"
	"net/url"
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
		body = CreateShortUrlResponse{
			Slug: createResult.Record.Slug,
			Host: c.Request.Host,
		}
	case enums.CreationResultAlreadyExists:
		status = http.StatusOK
		body = CreateShortUrlResponse{
			Slug: createResult.Record.Slug,
			Host: c.Request.Host,
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
	Slug string
	Host string
}

func (r CreateShortUrlResponse) MarshalJSON() ([]byte, error) {
	u := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   r.Slug,
	}

	return json.Marshal(struct {
		ShortUrl string `json:"short_url"`
	}{
		ShortUrl: u.String(),
	})
}
