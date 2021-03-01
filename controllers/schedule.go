package controllers

import (
	"chatbot/models"
	"fmt"
	"math/rand"
	"time"

	k "github.com/Alfex4936/kakao"
	"github.com/gin-gonic/gin"
)

// GetSchedule :POST /schedule, MUST: "cate": 카테고리 이름
// Carousel
func GetSchedule(c *gin.Context) {
	rand.Seed(time.Now().Unix()) // To pick a image for carousel card randomly

	var schedules []models.Schedule
	var length int = len(models.CardImages)
	carousel := k.Carousel{}.New(false, false)

	if _, err := dbmap.Select(&schedules, models.LoadSchedule); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
		return
	}

	for _, schedule := range schedules {
		card := k.BasicCard{}.New(false, false)
		card.Title = schedule.Content
		card.Desc = fmt.Sprintf("%v ~ %v", schedule.StartDate, schedule.EndDate)
		card.ThumbNail = k.ThumbNail{}.New(fmt.Sprintf("https://raw.githubusercontent.com/Alfex4936/kakaoChatbot-Ajou/main/imgs/%v.png", models.CardImages[rand.Int()%length]))
		carousel.Cards.Add(card)
	}

	c.PureJSON(200, carousel.Build())
}
