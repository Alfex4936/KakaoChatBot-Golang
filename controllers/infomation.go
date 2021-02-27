package controllers

import (
	"chatbot/models"
	"fmt"
	"strings"

	k "github.com/Alfex4936/kakao"
	"github.com/gin-gonic/gin"
)

const intel = "031-219-"

// AskWeather :POST /weather
// 메시지 종류: SimpleText
func AskWeather(c *gin.Context) {
	// 수원 영통구 현재 날씨 불러오기 (weather.naver.com)
	weather, err := models.GetWeather()
	if err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(err.Error(), nil))
		// http.StatusBadRequest 400을 보내고 싶으나, 400으로 하면 작동 X
		return
	}

	msg := fmt.Sprintf("📡 [수원시 영통구 날씨] 📡\n\n🌡 현재: %s, %s\n\n🌡 최저, 최고 온도: %s, %s\n\n☔ 낮, 밤 강수 확률: %s, %s\n\n😷 미세먼지: %s\n\n😷 초미세먼지: %s\n\n☀ 자외선: %s",
		weather.CurrentStatus, weather.CurrentTemp,
		weather.MinTemp, weather.MaxTemp,
		weather.RainDay, weather.RainNight,
		weather.FineDust, weather.UltraDust, weather.UV)

	c.PureJSON(200, k.SimpleText{}.Build(msg, nil))
}

// AskWeatherInCard :POST /weather2
// 메시지 종류: BasicCard
func AskWeatherInCard(c *gin.Context) {
	// 수원 영통구 현재 날씨 불러오기 (weather.naver.com)
	weather, err := models.GetWeather()
	if err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
		// http.StatusBadRequest 400을 보내고 싶으나, 400으로 하면 작동 X
		return
	}

	basicCard := k.BasicCard{}.New(true, true) // 썸네일, 버튼

	basicCard.Title = "[수원 영통구 기준]"

	msg := fmt.Sprintf("현재 날씨는 %s, %s\n최고기온 %s, 최저기온은 %s\n\n낮, 밤 강수 확률은 %s, %s\n미세먼지 농도는 %s\n자외선 지수는 %s",
		weather.CurrentStatus, weather.CurrentTemp,
		weather.MinTemp, weather.MaxTemp,
		weather.RainDay, weather.RainNight,
		weather.FineDust, weather.UV)

	basicCard.Desc = msg
	basicCard.Buttons.Add(k.LinkButton{}.New("자세히", models.NaverWeather))
	basicCard.ThumbNail = k.ThumbNail{FixedRatio: false}.New(fmt.Sprintf("https://raw.githubusercontent.com/Alfex4936/KakaoChatBot-Golang/main/imgs/%v.png?raw=true", weather.Icon))

	c.PureJSON(200, basicCard.Build())
}

// SearchProf :POST /prof, MUST: "keyword": 검색어
// 메시지 종류: CarouselCard
func SearchProf(c *gin.Context) {
	// JSON request parse
	var kjson k.Request
	if err := c.BindJSON(&kjson); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(err.Error(), nil))
		return
	}

	keyword := strings.TrimSpace(kjson.Action.Params["search"].(string))

	people, err := models.GetPeople(keyword)
	if err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(err.Error(), nil))
		return
	}

	if len(people.PhoneNumber) == 0 {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(fmt.Sprintf("%v 검색 결과가 없습니다.", keyword), nil))
		return
	} else if len(people.PhoneNumber) > 10 { // CarouselCard only supports 10 basicCards
		people.PhoneNumber = people.PhoneNumber[:10]
	}

	// Carousel (BasicCard)
	carousel := k.Carousel{}.New(false, false)

	for _, person := range people.PhoneNumber {
		// basicCard 케로셀에 담기
		card := k.BasicCard{}.New(false, true)
		card.Title = fmt.Sprintf("%v (%v)", person.Name, person.Email)
		card.Desc = fmt.Sprintf("전화번호: %v\n부서명: %v", intel+person.TelNo, person.DeptNm)

		// 전화 버튼, 웹 링크 버튼 케로셀에 담기
		card.Buttons.Add(k.CallButton{}.New("전화", intel+person.TelNo))
		card.Buttons.Add(k.LinkButton{}.New("이메일", fmt.Sprintf("mailto:%s?subject=안녕하세요.", person.Email)))

		carousel.Cards.Add(card)
	}

	c.PureJSON(200, carousel.Build())
}
