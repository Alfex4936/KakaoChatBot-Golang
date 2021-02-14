package controllers

import (
	"fmt"
	"kakao/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // it is a blank
)

/*
type Weather struct {
	MaxTemp       string // 최고 온도
	MinTemp       string // 최저 온도
	CurrentTemp   string // 현재 온도
	CurrentStatus string // 흐림, 맑음 ...
	RainDay       string // 강수 확률 낮
	RainNight     string // 강수 확률 밤
	FineDust      string // 미세먼지 [보통, 나쁨]
	UltraDust     string // 초미세먼지 [보통, 나쁨]
	UV            string // 자외선 지수 [낮음, ]
}
*/

// AskWeather :POST /weather
func AskWeather(c *gin.Context) {
	// 수원 영통구 현재 날씨 불러오기 (weather.naver.com)
	weather, _ := models.GetWeather()

	// Make a basicCard
	// template := gin.H{"outputs": []gin.H{{
	// 	"basicCard": gin.H{
	// 		"buttons":   []gin.H{{"action": "webLink", "label": "날씨 홈페이지 열기", "webLinkUrl": models.NaverWeather}},
	// 		"thumbnail": gin.H{"imageUrl": ""},
	// 		"title":     fmt.Sprintf("현재 수원 영통구 날씨: %s", weather.CurrentTemp),
	// 		"description": fmt.Sprintf("(해)<br>현재 %s<br>최저, 최고 온도: %s, %s<br>낮, 밤 강수 확률: %s, %s<br>미세먼지: %s<br>초미세먼지: %s<br>자외선: %s",
	// 			weather.CurrentStatus,
	// 			weather.MinTemp, weather.MaxTemp,
	// 			weather.RainDay, weather.RainNight,
	// 			weather.FineDust, weather.UltraDust, weather.UV),
	// 	},
	// }}}
	// basicCard := gin.H{"version": "2.0", "template": template}

	simpleText := models.BuildSimpleText(fmt.Sprintf("📡 [수원시 영통구 날씨] 📡\n\n🌡 현재: %s, %s\n\n🌡 최저, 최고 온도: %s, %s\n\n🌂 낮, 밤 강수 확률: %s, %s\n\n😷 미세먼지: %s\n\n😷 초미세먼지: %s\n\n☀ 자외선: %s",
		weather.CurrentStatus, weather.CurrentTemp,
		weather.MinTemp, weather.MaxTemp,
		weather.RainDay, weather.RainNight,
		weather.FineDust, weather.UltraDust, weather.UV))

	c.PureJSON(200, simpleText)
}
