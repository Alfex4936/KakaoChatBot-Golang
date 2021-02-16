package mappings

import (
	"kakao/controllers"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// Welcome message to check if server is running well.
func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{"welcome": "server is running well."})
}

// Router ...
var Router *gin.Engine

// CreateURLMappings to make endpoints
func CreateURLMappings() {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	Router = gin.New()
	// Apply the middleware to the router (works with groups too)
	Router.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "*",
		RequestHeaders: "*",
		Credentials:    true,
	}))

	// Router.Use(controllers.Cors())
	v1 := Router.Group("/v1")
	{
		v1.GET("/", Welcome)
		v1.GET("/notices/:num", controllers.GetAllNotices)
		v1.POST("/last", controllers.GetLastNotice)
		v1.POST("/today", controllers.GetTodayNotices)
		v1.POST("/yesterday", controllers.GetYesterdayNotices)
		v1.POST("/ask", controllers.AskCategory)
		v1.POST("/ask/category", controllers.ShowCategory)
		v1.POST("/schedule", controllers.GetSchedule)
		v1.POST("/search", controllers.SearchKeyword)
		// Infomation
		v1.POST("/info/weather", controllers.AskWeather)
		v1.POST("/info/prof", controllers.SearchProf)
	}
}
