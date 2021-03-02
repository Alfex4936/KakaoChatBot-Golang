package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// AjouMeal ...
var AjouMeal = os.Getenv("AJOU_MEAL")

// Meal ...
// 63 기숙사, 220 학생, 221 교직원
type Meal struct {
	IsSuccess string `json:"msgCode"`
	Data      struct {
		Breakfast string `json:"breakfast"`    // 아침
		Lunch     string `json:"lunch"`        // 점심
		Dinner    string `json:"dinner"`       // 저격
		SnackBar  string `json:"snackBar"`     // 분식
		Date      string `json:"menuDt"`       // 날짜
		Name      string `json:"restaurantNm"` // 식당 이름 (교직원식당(생활관 2층))
	} `json:"p018Text"`
}

// GetMeal ...
func GetMeal(place, when string) (Meal, error) {
	// when = "20210X0X"
	var meal Meal

	// POST
	jsonValue, _ := json.Marshal(map[string]string{"categoryId": place, "yyyymmdd": when})

	// Disable SSL authentication and post jsonValue
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post((AjouMeal), "application/json;charset=UTF-8", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
		return meal, err
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &meal)
	if err != nil {
		fmt.Println(err)
		return meal, err
	}

	if reflect.ValueOf(meal.Data).IsZero() {
		meal.IsSuccess = "empty"
	} else {
		// For JSON, replace all <br> into \n
		meal.Data.Breakfast = strings.ReplaceAll(meal.Data.Breakfast, "<br>", "\n")
		meal.Data.Lunch = strings.ReplaceAll(meal.Data.Lunch, "<br>", "\n")
		meal.Data.Dinner = strings.ReplaceAll(meal.Data.Dinner, "<br>", "\n")
		meal.Data.SnackBar = strings.ReplaceAll(meal.Data.SnackBar, "<br>", "\n")
	}

	return meal, nil
}
