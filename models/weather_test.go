package models

import (
	"strings"
	"testing"

	"github.com/anaskhan96/soup"
)

func TestGetWeatherIcon(t *testing.T) {
	resp, err := soup.Get("https://weather.naver.com/today/02117530?cpName=ACCUWEATHER")
	if err != nil {
		t.Errorf("Check your HTML connection.")
		return
	}
	doc := soup.HTMLParse(resp)
	img := doc.Find("div", "class", "today_weather").Find("i").Attrs()["data-ico"]

	expected := "ico_animation_wt"
	if strings.Contains(img, expected) {
		t.Logf("Found %q", img)
	} else {
		t.Errorf("Couldn't find!")
	}
}
