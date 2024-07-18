package handlers

import (
	"github.com/3XBAT/time-tracker/docs"
	"github.com/3XBAT/time-tracker/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/users", h.getUsers)
	router.GET("/users/:id", h.getUserByID)
	router.POST("/users", h.createUser)
	router.PUT("/users/:id", h.updateUser)
	router.DELETE("/users/:id", h.deleteUser)

	router.POST("/tasks/", h.createTask)
	router.PUT("/tasks/:id", h.updateTask)    //
	router.DELETE("/tasks/:id", h.deleteTask) //
	router.GET("tasks/", h.getTasks)
	return router
}
