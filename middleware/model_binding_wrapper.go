package middleware

import (
	"errors"
	"net/http"
	"url-shortener/controllers"
	"url-shortener/e"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ModelBindingWrapper[T any](controller controllers.Controller[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request T

		if err := c.ShouldBind(&request); err != nil {
			var verr validator.ValidationErrors
			if errors.As(err, &verr) {
				c.JSON(http.StatusBadRequest, e.ErrorResponse{
					Errors: e.FormatErrors(verr),
				})
			} else {
				c.Writer.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		controller.HandleRequest(c, request)
	}
}
