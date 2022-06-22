package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/anaskhan96/soup"
)

const naverWeather = "https://m.search.naver.com/search.naver?sm=tab_hty.top&where=nexearch&query=%EB%82%A0%EC%94%A8+%EB%A7%A4%ED%83%843%EB%8F%99&oquery=%EB%82%A0%EC%94%A8"
const naverWeatherIcon = "https://weather.naver.com/today/02117530?cpName=ACCUWEATHER"

// Weather 해외 기상은 일출부터 일몰 전이 낮, 일몰부터 일출 전이 밤
type weather struct {
	MaxTemp       string `json:"max_temp"`     // 최고 온도
	MinTemp       string `json:"min_temp"`     // 최저 온도
	CurrentTemp   string `json:"current_temp"` // 현재 온도
	CurrentStatus string `json:"current_stat"` // 흐림, 맑음 ...
	RainDay       string `json:"rain_day"`     // 강수 확률 낮
	RainNight     string `json:"rain_night"`   // 강수 확률 밤
	FineDust      string `json:"fine_dust"`    // 미세먼지
	UltraDust     string `json:"ultra_dust"`   // 초미세먼지
	UV            string `json:"uv"`           // 자외선 지수
	Sunset        string `json:"sunset"`       // 일몰
	Icon          string `json:"icon"`         // 날씨 아이콘 (ico_animation_wt?)
}

func TestGetWeatherIcon(t *testing.T) {
	// -short 쓰면 테스트 스킵
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }

	var weather weather

	resp, err := soup.Get(naverWeather)
	if err != nil {
		fmt.Println("Check your HTML connection.")
		return
	}

	resp2, err := soup.Get(naverWeatherIcon)
	if err != nil {
		fmt.Println("Check your HTML connection.")
		return
	}
	doc := soup.HTMLParse(resp)
	doc2 := soup.HTMLParse(resp2)

	currentTemp := doc.Find("div", "class", "temperature_text").Find("strong").Text() + "도"

	// ! 해외 기상은 일출부터 일몰 전이 낮, 일몰부터 일출 전이 밤
	maxTemp := doc.Find("span", "class", "highest")
	dayTemp := maxTemp.Text() + "도"
	dayTemp = strings.Replace(dayTemp, "°", "", 1)

	minTemp := doc.Find("span", "class", "lowest")
	nightTemp := minTemp.Text() + "도"
	nightTemp = strings.Replace(nightTemp, "°", "", 1)

	currentStatElem := doc.Find("span", "class", "weather")
	currentStat := currentStatElem.Text()

	// [미세먼지, 초미세먼지, 자외선, 일몰 시간]
	statuses := doc.FindAll("li", "class", "item_today")
	fineDust := statuses[0].Find("a").Find("span").Text()
	ultraDust := statuses[1].Find("a").Find("span").Text()
	UV := statuses[2].Find("a").Find("span").Text()
	sunset := statuses[3].Find("a").Find("span").Text()

	// 강우량
	rainElems := doc.FindAll("span", "class", "rainfall")
	dayRain := rainElems[0].Text()
	nightRain := rainElems[1].Text()

	// 날씨 아이콘
	i := doc2.Find("div", "class", "summary_img").Find("i").Attrs()
	img := i["data-ico"]

	if strings.Contains(i["class"], "night") {
		img += "_night"
	}

	// struct 값 변경
	weather.CurrentTemp = currentTemp
	weather.CurrentStatus = currentStat

	weather.MaxTemp = dayTemp // Assert that (day temp > night temp) in general
	weather.MinTemp = nightTemp

	weather.RainDay = dayRain
	weather.RainNight = nightRain

	weather.FineDust = fineDust
	weather.UltraDust = ultraDust
	weather.UV = UV
	weather.Sunset = sunset
	weather.Icon = fmt.Sprintf("https://raw.githubusercontent.com/Alfex4936/KakaoChatBot-Golang/main/imgs/%s.png?raw=true", img)

	fmt.Printf("%+v", weather)
}
