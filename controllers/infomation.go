package controllers

import (
	"chatbot/models"
	"fmt"
	"strings"
	"time"

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
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
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

// GetSeatsAvailable :POST /library
// 메시지 종류: BasicCard
func GetSeatsAvailable(c *gin.Context) {
	library, err := models.GetLibraryAvailable()
	if err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
		// http.StatusBadRequest 400을 보내고 싶으나, 400으로 하면 작동 X
		return
	}

	msg := []string{} // BasicCard.Description

	basicCard := k.BasicCard{}.New(false, true) // 썸네일, 버튼

	basicCard.Title = "[중앙도서관]"

	for _, lib := range library.Data.List { // 모든 도서관
		msg = append(msg, fmt.Sprintf("%v: %v/%v (잔여/전체)", lib.Name, lib.Available, lib.ActiveTotal))
	}

	basicCard.Desc = strings.Join(msg, "\n")
	basicCard.Buttons.Add(k.LinkButton{}.New("중앙도서관 홈페이지", "https://library.ajou.ac.kr/#/"))

	c.PureJSON(200, basicCard.Build())
}

// AskMeal :POST /meal
// 메시지 종류: SimpleText
func AskMeal(c *gin.Context) {
	// JSON request parse
	var kjson k.Request
	if err := c.BindJSON(&kjson); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil)) // http.StatusBadRequest
		return
	}

	userParams := kjson.Action.Params
	var postDate, postPlace string

	now := time.Now()
	nowStr := fmt.Sprintf("%v%02v%02v", now.Year(), int(now.Month()), now.Day())
	tomorrow := time.Now().Add(24 * time.Hour)
	tomorrowStr := fmt.Sprintf("%v%02v%02v", tomorrow.Year(), int(tomorrow.Month()), tomorrow.Day())

	when := userParams["when"].(string)
	place := userParams["place"].(string)

	switch when {
	case "오늘":
		postDate = nowStr
	case "내일":
		postDate = tomorrowStr
	default:
		postDate = nowStr
	}

	if strings.Contains(place, "기숙사") {
		place = "기숙사"
		postPlace = "63"
	} else if strings.Contains(place, "학생") || strings.Contains(place, "학식") {
		place = "학생"
		postPlace = "220"
	} else if strings.Contains(place, "교직원") {
		place = "교직원"
		postPlace = "221"
	} else {
		place = "학생"
		postPlace = "220"
	}

	// fmt.Println(postPlace, postDate)
	meal, err := models.GetMeal(postPlace, postDate)
	if err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
		// http.StatusBadRequest 400을 보내고 싶으나, 400으로 하면 작동 X
		return
	}

	var msg, breakfast, lunch, dinner, snack string

	if meal.Data.Breakfast != "" {
		breakfast = fmt.Sprintf("\n\n%v", meal.Data.Breakfast)
	}
	if meal.Data.Lunch != "" {
		lunch = fmt.Sprintf("\n\n%v", meal.Data.Lunch)
	}
	if meal.Data.Dinner != "" {
		dinner = fmt.Sprintf("\n\n%v", meal.Data.Dinner)
	}
	if meal.Data.SnackBar != "" {
		snack = fmt.Sprintf("\n\n%v", meal.Data.SnackBar)
	}

	if meal.IsSuccess == "empty" {
		msg = fmt.Sprintf("%s의 %s 식단이 없습니다!", when, place)
	} else {
		msg = fmt.Sprintf("[%v] %v%v%v%v%v",
			meal.Data.Date, meal.Data.Name,
			breakfast, lunch, dinner, snack)
	}

	var quickReplies k.Kakao
	quickReplies.Add(k.QuickReply{}.New("오늘 학식", "오늘 학식"))
	quickReplies.Add(k.QuickReply{}.New("오늘 기숙사", "오늘 기숙사"))
	quickReplies.Add(k.QuickReply{}.New("내일 교직원", "내일 교직원"))

	c.PureJSON(200, k.SimpleText{}.Build(msg, quickReplies))
}

// AskJob :POST /job
// 메시지 종류: CarouselCard
func AskJob(c *gin.Context) {
	job, err := models.GetJobsAvailable()
	if err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(err.Error(), nil))
		return
	}

	// Carousel (BasicCard)
	carousel := k.Carousel{}.New(false, false)

	for _, jobInfo := range job.Data {
		// basicCard 케로셀에 담기
		card := k.BasicCard{}.New(false, true)
		card.Title = jobInfo.Title
		card.Desc = fmt.Sprintf("공고 날짜: %v", jobInfo.Date)

		// 웹 링크 버튼 케로셀에 담기
		card.Buttons.Add(k.LinkButton{}.New("자세히", jobInfo.URL))

		carousel.Cards.Add(card)
	}

	c.PureJSON(200, carousel.Build())
}
