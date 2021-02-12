package controllers

import (
	"fmt"
	"kakao/models"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // it is a blank
)

// GetSchedule :POST /schedule, MUST: "cate": 카테고리 이름
func GetSchedule(c *gin.Context) {
	if err := models.CheckConnection(); err == false {
		errorMsg := models.SimpleText{Version: "2.0"}
		errorMsg.Template.Outputs.SimpleText.Text = "인터넷 연결 확인 후 잠시 후 시도하세요."
		c.JSON(404, errorMsg)
		return
	}
	rand.Seed(time.Now().Unix()) // To pick a image for carousel card randomly

	var schedules []models.Schedule
	var cards []gin.H
	var length int = len(models.CardImages)

	if _, err := dbmap.Select(&schedules, models.LoadSchedule); err != nil {
		errorMsg := models.SimpleText{Version: "2.0"}
		errorMsg.Template.Outputs.SimpleText.Text = err.Error()
		c.JSON(404, errorMsg)
		return
	}

	for _, schedule := range schedules {
		carousel := gin.H{"title": schedule.Content, "description": fmt.Sprintf("%v ~ %v", schedule.StartDate, schedule.EndDate),
			"thumbnail": gin.H{"imageUrl": fmt.Sprintf("https://raw.githubusercontent.com/Alfex4936/kakaoChatbot-Ajou/main/imgs/%v.png", models.CardImages[rand.Int()%length])},
		}
		cards = append(cards, carousel)
	}

	template := gin.H{"outputs": []gin.H{{"carousel": gin.H{"type": "basicCard", "items": cards}}}}
	carouselCard := gin.H{"version": "2.0", "template": template}

	c.PureJSON(200, carouselCard)
}
