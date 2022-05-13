package shorturls

import (
	"net/http"
	"url-shortener/e"
	"url-shortener/enums"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeleteShortUrlController struct {
	DB                    *gorm.DB
	DeleteShortUrlService *services.DeleteShortUrlService
}

// DeleteShortUrl  godoc
// @Summary      Delete an existing short URL
// @Description  Delete an existing short URL by supplying the slug.
// @Tags         shorturls
// @Accept       json
// @Produce      json
// @Param        slug  path  string  true  "slug of short URL to delete"
// @Success      204
// @Failure      404  {object}  e.ErrorResponse
// @Failure      500
// @Router       /shorturls/{slug} [delete]
func (controller *DeleteShortUrlController) HandleRequest(c *gin.Context) {
	slug := c.Param("slug")
	result := controller.DeleteShortUrlService.Delete(slug)

	if result.Error != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch result.Status {
	case enums.DeleteResultSuccessful:
		c.Writer.WriteHeader(http.StatusNoContent)
	case enums.DeleteResultNotFound:
		c.JSON(http.StatusNotFound, e.ErrorResponse{
			Errors: []e.ValidationError{
				{
					Field:  "Slug",
					Reason: "not found",
				},
			},
		})
	default:
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *DeleteShortUrlController) Register(r *gin.Engine) {
	r.DELETE("/api/v1/shorturls/:slug", controller.HandleRequest)
}
