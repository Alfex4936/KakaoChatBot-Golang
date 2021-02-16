package controllers

import (
	"fmt"
	"kakao/models"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // it is a blank
)

// AskCategory :POST /ask, MUST: "cate": 카테고리 이름
// 메시지 종류: SimpleText
func AskCategory(c *gin.Context) {
	categories := []string{"학사", "학사일정", "비교과",
		"장학", "취업", "사무",
		"행사", "파란학기제", "학술",
		"입학", "기타",
	}

	var replies []gin.H
	for _, cate := range categories {
		replies = append(replies, gin.H{"messageText": cate, "action": "message", "label": cate})
	}

	// Make a template
	template := gin.H{"outputs": []gin.H{{"simpleText": gin.H{"text": "무슨 공지를 보고 싶으신가요?"}}}}
	template["quickReplies"] = replies // Optional
	simpleText := gin.H{"version": "2.0", "template": template}

	c.PureJSON(200, simpleText)
}

// ShowCategory :POST /ask/category, MUST: "cate": 카테고리 이름
// 메시지 종류: SimpleText | ListCard
func ShowCategory(c *gin.Context) {
	// JSON request parse
	var kjson models.KakaoJSON
	if err := c.BindJSON(&kjson); err != nil {
		c.AbortWithStatusJSON(200, models.BuildSimpleText(err.Error())) // http.StatusBadRequest
		return
	}

	categories := map[string]int{
		"학사":    1,
		"비교과":   2,
		"장학":    3,
		"학술":    4,
		"입학":    5,
		"취업":    6,
		"사무":    7,
		"기타":    8,
		"행사":    166,
		"파란학기제": 167,
		"파란학기":  167,
		"학사일정":  168,
	}

	// Cast to string as cate parameter is an interface
	userCategory := strings.Replace(kjson.Action.Params["cate"].(string), " ", "", 1)
	url := fmt.Sprintf("%v?mode=list&srCategoryId=%v&srSearchKey=&srSearchVal=&articleLimit=5&article.offset=0", models.AjouLink, categories[userCategory])

	var notices []models.Notice = models.Parse(url, 5)
	if len(notices) == 0 {
		c.AbortWithStatusJSON(200, models.BuildSimpleText("아주대학교 홈페이지 서버 반응이 늦고 있네요. 잠시 후 다시 시도해보세요."))
		return
	}

	// Card
	items := []gin.H{}
	buttons := []gin.H{}
	header := gin.H{"title": fmt.Sprintf("%v 공지", userCategory)}

	// Add one care item
	buttons = append(buttons, gin.H{"label": "공유하기", "action": "share"})
	buttons = append(buttons, gin.H{"label": userCategory, "action": "webLink",
		"webLinkUrl": fmt.Sprintf("%v?mode=list&srCategoryId=%v", models.AjouLink, categories[userCategory])})

	for _, notice := range notices {
		if utf8.RuneCountInString(notice.Title) > 35 { // To count korean letters length correctly
			notice.Title = string([]rune(notice.Title)[0:32]) + "..."
		}
		description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])
		noticeJSON := gin.H{"title": notice.Title, "description": description, "link": gin.H{"web": notice.Link}}
		items = append(items, noticeJSON)
	}
	// Make a template
	template := gin.H{"outputs": []gin.H{{"listCard": gin.H{"header": header, "items": items, "buttons": buttons}}}}
	listCard := gin.H{"version": "2.0", "template": template}

	c.JSON(200, listCard)
}
