package models

import (
	"testing"
)

// go test -v ./models
func TestGetMeal(t *testing.T) {
	date := "20210308"
	meal, err := GetMeal("221", date)
	if err != nil {
		t.Errorf("Check your HTML connection.")
		return
	}

	if meal.IsSuccess == "empty" {
		t.Logf("No Meal")
	} else {
		t.Logf("%s %v", date, meal.Data)
	}
}
