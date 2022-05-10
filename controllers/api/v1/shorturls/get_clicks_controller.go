package shorturls

import (
	"net/http"
	"url-shortener/e"
	"url-shortener/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetShortUrlClicksController struct {
	DB *gorm.DB
}

type GetShortUrlClicksRequest struct {
	TimePeriod string `form:"time_period" binding:"oneof=24_HOURS 1_WEEK ALL_TIME,required"`
}

type GetShortUrlClicksResponse struct {
	Count      int64  `json:"count"`
	TimePeriod string `json:"time_period"`
}

func (controller *GetShortUrlClicksController) HandleRequest(c *gin.Context, request GetShortUrlClicksRequest) {
	slug := c.Param("slug")

	var count int64

	// The GROUP BY helps us get 0 results if the short_urls row does not exist.
	// We can check `RowsAffected` to understand if the query returned zero rows
	// (short url does not exist) or 1 row (short url does exist)
	countResult := controller.DB.
		Raw(`
			SELECT COUNT(*)
			FROM
				short_urls
				LEFT OUTER JOIN clicks ON clicks.short_url_id = short_urls.id
			WHERE
				short_urls.slug = ?
			GROUP BY short_urls.id
		`, slug).Scan(&count)

	if countResult.Error != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if countResult.RowsAffected == 1 {
		c.JSON(http.StatusOK, &GetShortUrlClicksResponse{
			Count:      count,
			TimePeriod: request.TimePeriod,
		})

		return
	}

	c.JSON(http.StatusNotFound, e.ErrorResponse{
		Errors: []e.ValidationError{
			e.ValidationError{
				Field:  "Slug",
				Reason: "not found",
			},
		},
	})
}

func (controller *GetShortUrlClicksController) Register(r *gin.Engine) {
	r.GET("/api/v1/shorturls/:slug/clicks", middleware.ModelBindingWrapper[GetShortUrlClicksRequest](controller))
}
