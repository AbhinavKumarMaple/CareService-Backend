package routes

import (
	scheduleController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/schedule"
	"github.com/gin-gonic/gin"
)

func ScheduleRoutes(router *gin.RouterGroup, controller scheduleController.IScheduleController) {
	scheduleRouter := router.Group("/schedules")
	{
		scheduleRouter.GET("/", controller.GetSchedules)
		scheduleRouter.POST("/", controller.CreateSchedule)
		scheduleRouter.GET("/today", controller.GetTodaySchedules)
		scheduleRouter.GET("/today/:assignedUserID", controller.GetTodaySchedulesByAssignedUserID)
		scheduleRouter.GET("/:id", controller.GetScheduleByID)
		scheduleRouter.PUT("/:id", controller.UpdateSchedule)
		scheduleRouter.POST("/:id/start", controller.StartSchedule)
		scheduleRouter.POST("/:id/end", controller.EndSchedule)
	}

	taskRouter := router.Group("/tasks")
	{
		taskRouter.POST("/:taskId/update", controller.UpdateTask)
	}
}