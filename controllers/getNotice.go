package controllers

import (
	"chatbot/models"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	k "github.com/Alfex4936/kakao"
	"github.com/gin-gonic/gin"
)

// GetAllNotices from my db :GET /notices/:num
func GetAllNotices(c *gin.Context) {
	quantity := c.Params.ByName("num")
	num, _ := strconv.ParseInt(quantity, 10, 64)

	var notices []models.Notice
	if _, err := dbmap.Select(&notices, models.PrintNotices, num); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
		return
	}

	c.PureJSON(200, notices)
}

// GetLastNotice :POST /last
// 메시지 종류: SimpleText | ListCard
func GetLastNotice(c *gin.Context) {
	// if err := dbmap.SelectOne(&notice, models.PrintNotices, 1); err != nil {
	// 	c.JSON(404, models.BuildSimpleText(err.Error()))
	// 	return
	// }
	var notice models.Notice = models.Parse("", 1)[0]

	// ListCard
	listCard := k.ListCard{}.New(false)
	listCard.Title = fmt.Sprintf("%v 공지", notice.Date)

	// Add button
	listCard.Buttons.Add(k.ShareButton{}.New("공유하기"))

	if utf8.RuneCountInString(notice.Title) > 35 { // To count korean letters length correctly
		notice.Title = string([]rune(notice.Title)[0:32]) + "..."
	}
	description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])

	listCard.Items.Add(k.ListItemLink{}.New(notice.Title, description, "", notice.Link))

	c.PureJSON(200, listCard.Build())
}

// GetTodayNotices :POST /today
// 메시지 종류: SimpleText | ListCard
func GetTodayNotices(c *gin.Context) {
	var notices []models.Notice = models.Parse("", 30)
	var label string

	now := time.Now()
	nowStr := fmt.Sprintf("%v.%02v.%02v", now.Year()%100, int(now.Month()), now.Day())
	// nowStr := "21.02.10"

	// Filtering out today's notice(s)
	for i, notice := range notices {
		if notice.Date != nowStr {
			notices = notices[:i]
			break
		}
	}

	// ListCard + QuickReplies
	listCard := k.ListCard{}.New(true)
	listCard.Title = fmt.Sprintf("%v) 오늘 공지", nowStr)

	// Add buttons
	listCard.Buttons.Add(k.ShareButton{}.New("공유하기"))

	if len(notices) > 5 {
		label = fmt.Sprintf("%v개 더보기", len(notices)-5)
		listCard.Buttons.Add(k.MsgButton{}.New(fmt.Sprintf("%v개 더보기", len(notices)-5), "더보기"))
		notices = notices[:5]
	} else {
		listCard.Buttons.Add(k.LinkButton{}.New(label, models.AjouLink))
		label = "아주대학교 공지"
	}

	if len(notices) == 0 {
		listCard.Items.Add(k.ListItem{}.New("공지가 없습니다!", "", "http://k.kakaocdn.net/dn/APR96/btqqH7zLanY/kD5mIPX7TdD2NAxgP29cC0/1x1.jpg"))
	} else {
		for _, notice := range notices {
			if utf8.RuneCountInString(notice.Title) > 35 { // To count korean letters length correctly
				notice.Title = string([]rune(notice.Title)[0:32]) + "..."
			}
			description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])

			listCard.Items.Add(k.ListItemLink{}.New(notice.Title, description, "", notice.Link))
		}
	}

	// Add Two quick replies
	listCard.QuickReplies.Add(k.QuickReply{}.New("오늘", "오늘 공지 보여줘"))
	listCard.QuickReplies.Add(k.QuickReply{}.New("어제", "어제 공지 보여줘"))

	c.PureJSON(200, listCard.Build())
}

// GetTodayMoreNotices Carousel + BasicCard + CarouselHeader
func GetTodayMoreNotices(c *gin.Context) {
	var notices []models.Notice = models.Parse("", 30)

	now := time.Now()
	today := fmt.Sprintf("%v월 %v일)", int(now.Month()), now.Day())
	nowStr := fmt.Sprintf("%v.%02v.%02v", now.Year()%100, int(now.Month()), now.Day())
	// nowStr := "21.02.10"

	// Filtering out today's notice(s)
	for i, notice := range notices {
		if notice.Date != nowStr {
			notices = notices[:i]
			break
		}
	}

	if len(notices) <= 5 { // 잘못된 접근 방지
		var quickReplies k.Kakao
		quickReplies.Add(k.QuickReply{}.New("오늘 공지", "ㅗ"))
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("공지가 5개 정도 뿐이에요!", quickReplies))
		return
	}

	// Carousel + QuickReplies
	carousel := k.Carousel{}.New(false, true)
	carousel.Header = k.CarouselHeader{}.New(fmt.Sprintf("오늘 공지 총 %d개", len(notices)), fmt.Sprintf("%d개를 더 불러왔습니다!", len(notices)-5), "http://k.kakaocdn.net/dn/APR96/btqqH7zLanY/kD5mIPX7TdD2NAxgP29cC0/1x1.jpg")

	if len(notices) == 0 {
		card := k.BasicCard{}.New(true, false)
		card.Title = "공지가 없습니다!"
		card.ThumbNail = k.ThumbNail{}.New("http://k.kakaocdn.net/dn/APR96/btqqH7zLanY/kD5mIPX7TdD2NAxgP29cC0/1x1.jpg")
	} else {
		for _, notice := range notices[5:] {
			// if utf8.RuneCountInString(notice.Title) > 35 { // To count korean letters length correctly
			// 	notice.Title = string([]rune(notice.Title)[0:33]) + ".."
			// }
			description := fmt.Sprintf("%v %v", today, notice.Writer)

			card := k.BasicCard{}.New(false, true)
			card.Title = description
			card.Desc = notice.Title
			card.Buttons.Add(k.LinkButton{}.New("공지 보기", notice.Link))
			carousel.Cards.Add(card)
		}
	}

	// Add Two quick replies
	// carousel.QuickReplies.Add(k.QuickReply{}.New("오늘", "오늘 공지 보여줘"))
	// carousel.QuickReplies.Add(k.QuickReply{}.New("어제", "어제 공지 보여줘"))

	c.PureJSON(200, carousel.Build())
}

// GetYesterdayNotices :POST /today
// 메시지 종류: SimpleText | ListCard
func GetYesterdayNotices(c *gin.Context) {
	yesterday := time.Now().Add(-24 * time.Hour)
	yesterdayStr := fmt.Sprintf("%v.%02v.%02v", yesterday.Year()%100, int(yesterday.Month()), yesterday.Day())

	var notices []models.Notice
	var label string

	if _, err := dbmap.Select(&notices, models.GetNoticesByDate, yesterdayStr); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("오류가 발생했습니다.\n:( 다시 시도해 주세요!", nil))
		return
	}

	// ListCard + QuickReplies
	listCard := k.ListCard{}.New(true)
	listCard.Title = fmt.Sprintf("%v) 어제 공지", yesterdayStr)

	// Add buttons
	listCard.Buttons.Add(k.ShareButton{}.New("공유하기"))

	if len(notices) > 5 {
		label = fmt.Sprintf("%v개 더보기", len(notices)-5)
		notices = notices[:5]
	} else {
		label = "아주대학교 공지"
	}
	listCard.Buttons.Add(k.LinkButton{}.New(label, models.AjouLink))

	// Python MakeJSON
	if len(notices) == 0 {
		listCard.Items.Add(k.ListItem{}.New("공지가 없습니다!", "", "http://k.kakaocdn.net/dn/APR96/btqqH7zLanY/kD5mIPX7TdD2NAxgP29cC0/1x1.jpg"))
	} else {
		for _, notice := range notices {
			if utf8.RuneCountInString(notice.Title) > 35 { // To count korean letters length correctly
				notice.Title = string([]rune(notice.Title)[0:32]) + "..."
			}
			description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])

			listCard.Items.Add(k.ListItemLink{}.New(notice.Title, description, "", notice.Link))
		}
	}

	// Add Two quick replies
	listCard.QuickReplies.Add(k.QuickReply{}.New("오늘", "오늘 공지 보여줘"))
	listCard.QuickReplies.Add(k.QuickReply{}.New("어제", "어제 공지 보여줘"))

	c.PureJSON(200, listCard.Build())
}
