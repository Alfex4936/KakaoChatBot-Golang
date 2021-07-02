package models

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
)

// ! TODO : 카카오톡 날씨에서나 현재에서 내일 날씨 peek

// const NaverWeather = "https://search.naver.com/search.naver?query=수원영통구날씨"
// NaverWeather ...
const NaverWeather = "https://weather.naver.com/today/02117530?cpName=ACCUWEATHER" // 아큐웨더 제공 날씨

// Weather 해외 기상은 일출부터 일몰 전이 낮, 일몰부터 일출 전이 밤
type Weather struct {
	MaxTemp       string // 최고 온도
	MinTemp       string // 최저 온도
	CurrentTemp   string // 현재 온도
	CurrentStatus string // 흐림, 맑음 ...
	RainDay       string // 강수 확률 낮
	RainNight     string // 강수 확률 밤
	FineDust      string // 미세먼지
	UltraDust     string // 초미세먼지
	UV            string // 자외선 지수
	Icon          string // 날씨 아이콘 (ico_animation_wt?)
}

// GetWeather is a function that parses suwon's weather today
func GetWeather() (Weather, error) {
	var weather Weather

	resp, err := soup.Get(NaverWeather)
	if err != nil {
		fmt.Println("Check your HTML connection.")
		return weather, err // nil
	}
	doc := soup.HTMLParse(resp)

	currentTempInt := doc.Find("strong", "class", "current").Text() + "도"

	currentStatus := doc.Find("span", "class", "weather")

	// ! 해외 기상은 일출부터 일몰 전이 낮, 일몰부터 일출 전이 밤
	temps := doc.FindAll("span", "class", "data")
	DayTemp := temps[0].Text() + "도"
	NightTemp := temps[1].Text() + "도"

	rains := doc.FindAll("span", "class", "rainfall")
	DayRain := rains[0].Text()
	NightRain := rains[1].Text()

	// [미세먼지, 초미세먼지, 자외선, 일몰 시간]
	statuses := doc.FindAll("em", "class", "level_text")

	// 날씨 아이콘
	i := doc.Find("div", "class", "today_weather").Find("i").Attrs()
	img := i["data-ico"]

	// struct 값 변경
	weather.CurrentTemp = currentTempInt
	weather.CurrentStatus = currentStatus.Text()

	weather.MaxTemp = DayTemp // Assert that (day temp > night temp) in general
	weather.MinTemp = NightTemp

	weather.RainDay = DayRain
	weather.RainNight = NightRain

	weather.FineDust = statuses[0].Text()
	weather.UltraDust = statuses[1].Text()
	weather.UV = statuses[2].Text()

	if strings.Contains(i["class"], "night") {
		weather.Icon = img + "_night"
	} else {
		weather.Icon = img
	}
	return weather, nil
}
