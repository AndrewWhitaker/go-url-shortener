package controllers

import "github.com/gin-gonic/gin"

type Controller[T any] interface {
	HandleRequest(c *gin.Context, model T)
}

type RegistrableController interface {
	Register(r *gin.Engine)
}
