package routes

import (
	"caregiver/src/infrastructure/rest/controllers/user"
	"caregiver/src/infrastructure/rest/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, controller user.IUserController) {
	u := router.Group("/user")
	u.Use(middlewares.AuthJWTMiddleware())
	{
		u.POST("/", controller.NewUser)
		u.GET("/", controller.GetAllUsers)
		u.GET("/:id", controller.GetUsersByID)
		u.PUT("/:id", controller.UpdateUser)
		u.DELETE("/:id", controller.DeleteUser)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
