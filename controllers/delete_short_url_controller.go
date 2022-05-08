package controllers

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
				e.ValidationError{
					Field:  "Slug",
					Reason: "not found",
				},
			},
		})
	default:
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
}
