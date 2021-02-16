package controllers

import (
	"fmt"
	"kakao/models"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // it is a blank
)

const intel = "031-219-"

// AskWeather :POST /weather
// 메시지 종류: SimpleText
func AskWeather(c *gin.Context) {
	// 수원 영통구 현재 날씨 불러오기 (weather.naver.com)
	weather, err := models.GetWeather()
	if err != nil {
		c.AbortWithStatusJSON(200, models.BuildSimpleText(err.Error())) // http.StatusBadRequest
		return
	}

	simpleText := models.BuildSimpleText(fmt.Sprintf("📡 [수원시 영통구 날씨] 📡\n\n🌡 현재: %s, %s\n\n🌡 최저, 최고 온도: %s, %s\n\n☔ 낮, 밤 강수 확률: %s, %s\n\n😷 미세먼지: %s\n\n😷 초미세먼지: %s\n\n☀ 자외선: %s",
		weather.CurrentStatus, weather.CurrentTemp,
		weather.MinTemp, weather.MaxTemp,
		weather.RainDay, weather.RainNight,
		weather.FineDust, weather.UltraDust, weather.UV))

	c.PureJSON(200, simpleText)
}

// SearchProf :POST /prof, MUST: "keyword": 검색어
// 메시지 종류: CarouselCard
func SearchProf(c *gin.Context) {
	// JSON request parse
	var kjson models.KakaoJSON
	if err := c.BindJSON(&kjson); err != nil {
		c.AbortWithStatusJSON(200, models.BuildSimpleText(err.Error()))
		return
	}

	keyword := strings.TrimSpace(kjson.Action.Params["search"].(string))

	people, err := models.GetPeople(keyword)
	if err != nil {
		c.AbortWithStatusJSON(200, models.BuildSimpleText(err.Error()))
		return
	}

	if len(people.PhoneNumber) == 0 {
		c.AbortWithStatusJSON(200, models.BuildSimpleText(fmt.Sprintf("%v 검색 결과가 없습니다.", keyword)))
		return
	} else if len(people.PhoneNumber) > 10 { // CarouselCard only supports 10 basicCards
		people.PhoneNumber = people.PhoneNumber[:10]
	}

	var cards []gin.H

	for _, person := range people.PhoneNumber {
		carousel := gin.H{"title": fmt.Sprintf("%v (%v)", person.Name, person.Email), "description": fmt.Sprintf("전화번호: %v\n부서명: %v", intel+person.TelNo, person.DeptNm),
			//"thumbnail": gin.H{"imageUrl": "https://raw.githubusercontent.com/Alfex4936/kakaoChatbot-Ajou/main/imgs/people.png"},
			"buttons": []gin.H{{"action": "phone", "label": "전화", "phoneNumber": intel + person.TelNo}, {"action": "webLink", "label": "이메일", "webLinkUrl": fmt.Sprintf("mailto:%s?subject=안녕하세요.", person.Email)}},
		}
		cards = append(cards, carousel)
	}

	template := gin.H{"outputs": []gin.H{{"carousel": gin.H{"type": "basicCard", "items": cards}}}}
	carouselCard := gin.H{"version": "2.0", "template": template}

	c.PureJSON(200, carouselCard)
}
