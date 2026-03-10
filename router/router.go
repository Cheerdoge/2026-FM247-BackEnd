package router

import (
	handler "2026-FM247-BackEnd/handlers"
	middleware "2026-FM247-BackEnd/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	authhandler *handler.AuthHandler,
	avatarHandler *handler.AvatarHandler,
	todohandler *handler.TodoHandler,
	studydatahandler *handler.StudyDataHandler,
	musichandler *handler.MusicHandler,
) {
	publicGroup := r.Group("/api")
	{
		// 用户相关
		publicGroup.POST("/auth/register", authhandler.RegisterUserHandler)
		publicGroup.POST("/auth/login", authhandler.LoginHandler)
	}

	authGroup := r.Group("/api")
	authGroup.Use(middleware.AuthMiddleware(authhandler.Tokenservice))
	{
		// 用户相关
		authGroup.POST("/auth/logout", authhandler.LogoutHandler)
		authGroup.POST("/auth/cancel", authhandler.CancelHandler)
		authGroup.POST("/user/update_info", authhandler.UpdateUserInfoHandler)
		authGroup.POST("/user/update_email", authhandler.UpdateEmailHandler)
		authGroup.POST("/user/update_password", authhandler.UpdatePasswordHandler)
		authGroup.GET("/user/info", authhandler.GetUserInfoHandler)

		authGroup.POST("/user/avatar", avatarHandler.UploadAvatar)

		// 待办事项相关
		authGroup.POST("/todos", todohandler.CreateTodo)
		authGroup.GET("/todos", todohandler.GetTodos)
		authGroup.GET("/todos/:id", todohandler.GetTodoByID)
		authGroup.PUT("/todos/:id", todohandler.UpdateTodo)
		authGroup.DELETE("/todos/:id", todohandler.DeleteTodo)

		// 学习数据相关
		authGroup.POST("/studydata", studydatahandler.AddStudyData)
		authGroup.GET("/studydata/daily", studydatahandler.GetDailyStudyData)
		authGroup.GET("/studydata/total", studydatahandler.GetTotalStudyData)
		authGroup.GET("/studydata/weekly", studydatahandler.GetWeekStudyData)
		authGroup.GET("/studydata/monthly", studydatahandler.GetMonthlyStudyData)
		authGroup.GET("/studydata/yearly", studydatahandler.GetYearStudyData)

		// 音乐相关
		authGroup.GET("/music", musichandler.GetAllMusic)
		authGroup.POST("/music", musichandler.UploadMusic)
	}
}
