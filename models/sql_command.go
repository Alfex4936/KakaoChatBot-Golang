package models

const (
	PrintNotices     = "SELECT * FROM ajou_notices ORDER BY id DESC LIMIT ?"
	GetNoticesByDate = "SELECT * FROM ajou_notices WHERE date = ? ORDER BY id DESC"
	LoadSchedule     = "SELECT * FROM ajou_sched LIMIT 10"
)
