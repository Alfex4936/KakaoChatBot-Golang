package controllers

import (
	"fmt"
	"kakao/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func checkConnection() (ok bool) {
	_, err := http.Get(models.AjouLink)
	if err != nil {
		return false
	}
	return true
}

func GetAllNotices(c *gin.Context) {
	quantity := c.Params.ByName("num")
	num, _ := strconv.ParseInt(quantity, 10, 64)

	var notices []models.Notice
	if _, err := dbmap.Select(&notices, models.PrintNotices, num); err != nil {
		errorMsg := models.SimpleText{Version: "2.0"}
		errorMsg.Template.Outputs.SimpleText.Text = err.Error()
		c.JSON(404, errorMsg)
		return
	}

	c.PureJSON(200, notices)
}

func GetLastNotice(c *gin.Context) {
	if err := checkConnection(); err != true {
		errorMsg := models.SimpleText{Version: "2.0"}
		errorMsg.Template.Outputs.SimpleText.Text = "인터넷 연결을 확인하세요."
		c.JSON(404, errorMsg)
		return
	}
	var notice models.Notice
	// c.Bind(&notice)
	if err := dbmap.SelectOne(&notice, models.PrintNotices, 1); err != nil {
		errorMsg := models.SimpleText{Version: "2.0"}
		errorMsg.Template.Outputs.SimpleText.Text = err.Error()
		c.JSON(404, errorMsg)
		return
	}
	c.PureJSON(200, notice)
}

func GetTodayNotices(c *gin.Context) {
	if err := checkConnection(); err != true {
		errorMsg := models.SimpleText{Version: "2.0"}
		errorMsg.Template.Outputs.SimpleText.Text = "인터넷 연결을 확인하세요."
		c.JSON(404, errorMsg)
		return
	}

	var notices []models.Notice = models.Parse(30)
	var label string

	now := time.Now()
	nowStr := fmt.Sprintf("%v.%02v.%v", now.Year()%100, int(now.Month()), now.Day())
	// nowStr := "21.02.10"

	// Filtering out today's notice(s)
	for i, notice := range notices {
		if notice.Date != nowStr {
			notices = notices[:i]
			break
		}
	}

	// Card
	items := []gin.H{}
	buttons := []gin.H{}
	header := gin.H{"title": fmt.Sprintf("%v) 오늘 공지", nowStr)}

	// Add one care item
	buttons = append(buttons, gin.H{"label": "공유하기", "action": "share"})

	if len(notices) > 5 {
		label = fmt.Sprintf("%v개 더보기", len(notices)-5)
	} else {
		label = "아주대학교 공지"
	}
	buttons = append(buttons, gin.H{"label": label, "action": "webLink", "webLinkUrl": models.AjouLink})

	if len(notices) == 0 {
		items = append(items, gin.H{"title": "공지가 없습니다!", "imageUrl": "http://k.kakaocdn.net/dn/APR96/btqqH7zLanY/kD5mIPX7TdD2NAxgP29cC0/1x1.jpg"})
	} else {
		for _, notice := range notices {
			noticeJSON := gin.H{"title": notice.Title, "description": notice.Writer, "link": gin.H{"web": notice.Link}}
			items = append(items, noticeJSON)
		}
	}

	// QuickReplies [Optional]
	quickReplies := []gin.H{}

	// Add Two quick replies
	quickReply1 := gin.H{"messageText": "오늘 공지 보여줘", "action": "message", "label": "오늘"}
	quickReply2 := gin.H{"messageText": "어제 공지 보여줘", "action": "message", "label": "어제"}
	quickReplies = append(quickReplies, quickReply1, quickReply2)

	// Make a template
	template := gin.H{"outputs": []gin.H{{"listCard": gin.H{"header": header, "items": items, "buttons": buttons}}}}
	template["quickReplies"] = quickReplies // Optional
	listCard := gin.H{"version": "2.0", "template": template}
	c.PureJSON(200, listCard)
}
