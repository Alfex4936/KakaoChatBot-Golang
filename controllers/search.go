package controllers

import (
	"fmt"
	"kakao/models"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // it is a blank
)

// SearchKeyword :POST /serach, MUST: "cate": 카테고리 이름
func SearchKeyword(c *gin.Context) {
	// JSON request parse
	var kjson models.KakaoJSON
	if err := c.BindJSON(&kjson); err != nil {
		c.JSON(200, models.BuildSimpleText(err.Error())) // http.StatusBadRequest
		return
	}

	// Cast to string as cate parameter is an interface
	userKeyword := kjson.Action.Params
	if _, ok := userKeyword["sys_text"]; !ok {
		// Make a template
		template := gin.H{"outputs": []gin.H{{"simpleText": gin.H{"text": "2021 검색과 같이 검색어를 같이 입력하세요."}}}}
		template["quickReplies"] = []gin.H{gin.H(models.BuildQuickReply("2021 검색", "2021 검색"))} // Optional
		simpleText := gin.H{"version": "2.0", "template": template}
		c.JSON(400, simpleText)
		return
	}
	keyword := userKeyword["sys_text"].(string)
	url := fmt.Sprintf("%v?mode=list&srSearchKey=&srSearchVal=%v&articleLimit=7&article.offset=0", models.AjouLink, url.QueryEscape(strings.TrimSpace(keyword)))

	var notices []models.Notice = models.Parse(url, 7)
	var replies []gin.H
	var label string

	if len(notices) == 0 {
		c.JSON(400, models.BuildSimpleText(fmt.Sprintf("%v에 관한 글이 없어요.", keyword)))
		return
	}
	// Card
	items := []gin.H{}
	buttons := []gin.H{}
	if utf8.RuneCountInString(keyword) > 12 {
		keyword = keyword[0:12]
	}
	header := gin.H{"title": fmt.Sprintf("%v 결과", keyword)}

	// Add one care item
	buttons = append(buttons, gin.H{"label": "공유하기", "action": "share"})
	if len(notices) > 5 {
		label = "더보기"
	} else {
		label = "홈페이지 보기"
	}
	buttons = append(buttons, gin.H{"label": label, "action": "webLink",
		"webLinkUrl": fmt.Sprintf("%v?mode=list&srSearchKey=&srSearchVal=%v", models.AjouLink, keyword)})

	// Python makeJSONwithDate(postId, postTitle, postDate, postLink, postWriter)
	for _, notice := range notices {
		description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])
		noticeJSON := gin.H{"title": notice.Title, "description": description, "link": gin.H{"web": notice.Link}}
		items = append(items, noticeJSON)
	}

	replies = append(replies, gin.H(models.BuildQuickReply("등록금 검색", "등록금 검색")))
	replies = append(replies, gin.H(models.BuildQuickReply("이벤트 검색", "이벤트 검색")))
	replies = append(replies, gin.H(models.BuildQuickReply("코로나 검색", "코로나 검색")))

	// Make a template
	template := gin.H{"outputs": []gin.H{{"listCard": gin.H{"header": header, "items": items, "buttons": buttons}}}}
	template["quickReplies"] = replies // Optional
	listCard := gin.H{"version": "2.0", "template": template}

	c.PureJSON(200, listCard)
}
