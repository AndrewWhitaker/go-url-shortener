package shorturls

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetShortUrlController struct {
	DB *gorm.DB
}

func (controller *GetShortUrlController) HandleRequest(c *gin.Context) {
	slug := c.Param("slug")

}
