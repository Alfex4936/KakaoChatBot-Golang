package mappings

import (
	"kakao/controllers"

	"github.com/gin-gonic/gin"
)

// Router ...
var Router *gin.Engine

// CreateURLMappings to make endpoints
func CreateURLMappings() {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	Router = gin.Default()

	Router.Use(controllers.Cors())
	v1 := Router.Group("/v1")
	{
		v1.GET("/notices/:num", controllers.GetAllNotices)
		v1.POST("/last/", controllers.GetLastNotice)
		v1.POST("/today/", controllers.GetTodayNotices)
		v1.POST("/ask/", controllers.AskCategory)
		v1.POST("/ask/category", controllers.ShowCategory)
	}
}
