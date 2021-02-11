package mappings

import (
	"kakao/controllers"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func CreateUrlMappings() {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	Router = gin.Default()

	Router.Use(controllers.Cors())
	v1 := Router.Group("/v1")
	{
		v1.GET("/notices/:num", controllers.GetAllNotices)
		v1.POST("/last/", controllers.GetLastNotice)
		v1.POST("/today/", controllers.GetTodayNotices)
	}
}
