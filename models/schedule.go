package models

// Schedule ...
type Schedule struct {
	ID        int64  `db:"id" json:"id"`
	Content   string `db:"content" json:"title"`
	StartDate string `db:"start_date" json:"date"`
	EndDate   string `db:"end_date" json:"link"`
}

// CardImages ...
var CardImages = []string{"ajou_carousel", "ajou_carousel_1", "ajou_carousel_2"}
