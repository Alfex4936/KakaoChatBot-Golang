package models

import (
	"strings"
	"testing"

	"github.com/anaskhan96/soup"
)

func TestGetWeatherIcon(t *testing.T) {
	// -short 쓰면 테스트 스킵
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }

	resp, err := soup.Get("https://weather.naver.com/today/02117530?cpName=ACCUWEATHER")
	if err != nil {
		t.Errorf("Check your HTML connection.")
		return
	}
	doc := soup.HTMLParse(resp)
	i := doc.Find("div", "class", "today_weather").Find("i")
	img := i.Attrs()["data-ico"]
	classes := i.Attrs()

	if strings.Contains(classes["class"], "night") {
		t.Logf("It's night!: %q", classes["class"])
	}

	// 밤에 night가 포함
	expected := "ico_animation_wt"
	if strings.Contains(img, expected) {
		t.Logf("Found %q", img)
		t.Logf("Found %q", classes)
	} else {
		t.Errorf("Couldn't find!")
	}
}
