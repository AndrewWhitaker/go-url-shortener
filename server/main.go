package server

import (
	"database/sql"
	"url-shortener/controllers"
	"url-shortener/controllers/api/v1/shorturls"
	"url-shortener/controllers/api/v1/shorturls/clicks"
	"url-shortener/db"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServerConfig struct {
	DB *sql.DB
}

func SetupServer(cfg *ServerConfig) *gin.Engine {
	db, err := db.ConnectDatabase(cfg.DB)

	if err != nil {
		panic("Failed to connect to database")
	}

	r := gin.Default()

	controllers := BuildControllers(db)

	for _, c := range controllers {
		c.Register(r)
	}

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
	}
}