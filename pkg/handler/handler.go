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
		auth.POST("/test", h.test)
		auth.GET("/get-user-fio", h.GetUserFIO)
		auth.POST("/get-user-not-access/:guid", h.GetUserNotAccess)
	}

	auth_start := router.Group("/auth-start")
	{
		auth_start.POST("/register", h.register)
		auth_start.POST("/login", h.login)
		auth_start.POST("/check-datatime-token", h.CheckToken)
		auth_start.POST("/create-reset-password", h.CreateResetPassword)
		auth_start.POST("/check-reset-password", h.CheckResetPassword)
		auth_start.POST("/use-reset-password", h.UserResetPassword)
	}

	role := router.Group("/role", h.userIdentity)
	{
		role.POST("/issue-access", h.IssueAccess)
		role.POST("/check-admin", h.CheackAdmin)
		role.POST("/check-access", h.CheackAccess)
		role.POST("/createInvite", h.createInvite)
		role.POST("/useInvite", h.useInvite)
	}

	plan := router.Group("/plan")
	{
		plan.GET("/get-mas-plan/:guid_program", h.GetMasPlans)
		plan.GET("/get-work-program/:guid_plan", h.GetWorkProgram)
		plan.GET("/get-field-plan/:guid_plan/:key_field", h.GetField)
		plan.POST("generate-word", h.GenerateWord)
		plan.GET("/get-name/:guid", h.GetNamePlans)
		plan.POST("/create-group-plans/:guid_faculty", h.createGroupPlans)
		plan.POST("/save-plan/:guid_plan/:key_field", h.SavePlan)
	}

	plan_auth := router.Group("/plan-auth", h.userIdentity)
	{
		plan_auth.POST("/create-group-plans/:guid_faculty", h.createGroupPlans)
		plan_auth.POST("/save-plan/:guid_plan/:key_field", h.SavePlan)
		plan_auth.POST("/copy-plan/:guid_to/:guid_from", h.CopyPlan)
	}

	faculty := router.Group("/faculty")
	{
		faculty.GET("/get-mas-faculty", h.GetMasFaculty)
		faculty.GET("/get-name/:guid", h.GetNameFaculte)
	}

	program := router.Group("/program")
	{
		program.GET("get-mas-program/:guid_faculty", h.GetMasProgram)
		program.GET("get-name-p-and-f/:guid", h.GetNameProgramAndFaculty)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://192.168.1.56:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}
