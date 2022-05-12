package clicks

import (
	"net/http"
	"url-shortener/e"
	"url-shortener/enums"
	"url-shortener/middleware"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
)

type GetShortUrlClicksController struct {
	GetClicksService *services.GetClicksService
}

type GetShortUrlClicksRequest struct {
	TimePeriod string `form:"time_period" binding:"oneof=24_HOURS 1_WEEK ALL_TIME,required"`
}

type GetShortUrlClicksResponse struct {
	Count      int64  `json:"count"`
	TimePeriod string `json:"time_period"`
}

// GetShortUrlClicks  godoc
// @Summary      Get clicks for a short URL
// @Description  Get clicks (statistics) for a short URL. Time periods of all time, 24 hours, and 1 week are permitted.
// @Tags         shorturls
// @Accept       json
// @Produce      json
// @Param        slug         path      string  true  "slug of short URL to retrieve statistics for"
// @Param        time_period  query     string  true  "time period to retrieve statistics for"  Enums(24_HOURS, 1_WEEK, ALL_TIME)
// @Success      200          {object}  GetShortUrlClicksResponse
// @Failure      404          {object}  e.ErrorResponse
// @Failure      500
// @Router       /shorturls/{slug}/clicks [get]
func (controller *GetShortUrlClicksController) HandleRequest(c *gin.Context, request GetShortUrlClicksRequest) {
	slug := c.Param("slug")

	result := controller.GetClicksService.GetClicks(
		slug, enums.GetClicksTimePeriodAllTime,
	)

	var status int
	var body interface{}

	if result.Error != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch result.Status {
	case enums.GetClicksResultSuccessful:
		status = http.StatusOK
		body = GetShortUrlClicksResponse{
			Count:      result.Count,
			TimePeriod: request.TimePeriod,
		}
	case enums.GetClicksResultNotFound:
		status = http.StatusNotFound
		body = e.ErrorResponse{
			Errors: []e.ValidationError{
				e.ValidationError{
					Field:  "Slug",
					Reason: "not found",
				},
			},
		}
	}

	c.JSON(status, body)
}

func (controller *GetShortUrlClicksController) Register(r *gin.Engine) {
	r.GET("/api/v1/shorturls/:slug/clicks", middleware.ModelBindingWrapper[GetShortUrlClicksRequest](controller))
}
