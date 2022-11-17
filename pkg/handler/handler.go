package handler

import (
	"time"

	"github.com/RucardTomsk/course_book/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth", h.userIdentity)
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/test", h.test)
	}

	plan := router.Group("/plan")
	{
		plan.POST("/create-group-plans/:guid_faculty", h.createGroupPlans)
		plan.GET("/get-mas-plan/:guid_program", h.GetMasPlans)
		plan.GET("/get-work-program/:guid_plan", h.GetWorkProgram)
		plan.POST("/save-plan/:guid_plan/:key_field", h.SavePlan)
		plan.GET("get-field-plan/:guid_plan/:key_field", h.GetField)
	}

	faculty := router.Group("/faculty")
	{
		faculty.GET("/get-mas-faculty", h.GetMasFaculty)
	}

	program := router.Group("/program")
	{
		program.GET("get-mas-program/:guid_faculty", h.GetMasProgram)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}
