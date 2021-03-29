package models

// Schedule ...
type Schedule struct {
	ID        int64  `db:"id" json:"id"`
	Content   string `db:"content" json:"content"`
	StartDate string `db:"start_date" json:"start_date"`
	EndDate   string `db:"end_date" json:"end_date"`
}

// CardImages ...
// 겨울 ajou_carousel_2
var CardImages = []string{"ajou_carousel", "ajou_carousel_1"}
