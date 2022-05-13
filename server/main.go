package server

import (
	"url-shortener/controllers"
	"url-shortener/controllers/api/v1/shorturls"
	"url-shortener/controllers/api/v1/shorturls/clicks"
	_ "url-shortener/docs"
	"url-shortener/services"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/gorm"
)

type ServerConfig struct {
	DB *gorm.DB
}

func SetupServer(cfg *ServerConfig) *gin.Engine {
	r := gin.Default()

	controllers := BuildControllers(cfg.DB)

	for _, c := range controllers {
		c.Register(r)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(static.Serve("/", static.LocalFile("./views", true)))

	return r
}

func BuildControllers(db *gorm.DB) []controllers.RegistrableController {
	createShortUrlService := &services.CreateShortUrlService{DB: db}
	deleteShortUrlService := &services.DeleteShortUrlService{DB: db}
	getClicksService := &services.GetClicksService{DB: db, Clock: services.SystemClock{}}

	createShortUrlController := shorturls.CreateShortUrlController{
		CreateShortUrlService: createShortUrlService,
	}

	deleteShortUrlController := shorturls.DeleteShortUrlController{
		DeleteShortUrlService: deleteShortUrlService,
	}

	getShortUrlController := shorturls.GetShortUrlController{
		DB: db,
	}

	listShortUrlsController := shorturls.ListShortUrlsController{
		DB: db,
	}

	getShortUrlClicksController := clicks.GetShortUrlClicksController{
		GetClicksService: getClicksService,
	}

	accessShortUrlController := controllers.AccessShortUrlController{
		DB: db,
	}

	return []controllers.RegistrableController{
		&createShortUrlController,
		&deleteShortUrlController,
		&accessShortUrlController,
		&getShortUrlClicksController,
		&getShortUrlController,
		&listShortUrlsController,
	}
}
