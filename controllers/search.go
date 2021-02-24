package controllers

import (
	"chatbot/models"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	k "github.com/Alfex4936/kakao"
	"github.com/gin-gonic/gin"
)

// SearchKeyword :POST /search, MUST: "cate": 카테고리 이름
func SearchKeyword(c *gin.Context) {
	// JSON request parse
	var kjson k.Request
	if err := c.BindJSON(&kjson); err != nil {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(err.Error(), nil)) // http.StatusBadRequest
		return
	}

	// Cast to string as cate parameter is an interface
	userKeyword := kjson.Action.Params
	if _, ok := userKeyword["sys_text"]; !ok {
		var quickReplies k.Kakao
		quickReplies.Add(k.QuickReply{}.New("2021 검색", "2021 검색"))
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build("2021 검색과 같이 검색어를 같이 입력하세요.", quickReplies))
		return
	}
	keyword := strings.TrimSpace(userKeyword["sys_text"].(string))
	url := fmt.Sprintf("%v?mode=list&srSearchKey=&srSearchVal=%v&articleLimit=7&article.offset=0", models.AjouLink, url.QueryEscape(keyword))

	var notices []models.Notice = models.Parse(url, 7)
	var label string

	if len(notices) == 0 {
		c.AbortWithStatusJSON(200, k.SimpleText{}.Build(fmt.Sprintf("%v에 관한 글이 없어요.", keyword), nil))
		return
	}

	// ListCard + QuickReplies
	listCard := k.ListCard{}.New(true)

	if utf8.RuneCountInString(keyword) > 12 {
		keyword = keyword[0:12]
	}
	listCard.Title = fmt.Sprintf("%v 결과", keyword)

	// Add one share button
	listCard.Buttons.Add(k.ShareButton{}.New("공유하기"))

	if len(notices) > 5 {
		label = "더보기"
	} else {
		label = "홈페이지 보기"
	}
	listCard.Buttons.Add(k.LinkButton{}.New(label, fmt.Sprintf("%v?mode=list&srSearchKey=&srSearchVal=%v", models.AjouLink, keyword)))

	// Python makeJSONwithDate(postId, postTitle, postDate, postLink, postWriter)
	for _, notice := range notices {
		description := fmt.Sprintf("%v %v", notice.Writer, notice.Date[len(notice.Date)-5:])
		listCard.Items.Add(k.ListItemLink{}.New(notice.Title, description, "", notice.Link))
	}

	listCard.QuickReplies.Add(k.QuickReply{}.New("등록금 검색", "등록금 검색"))
	listCard.QuickReplies.Add(k.QuickReply{}.New("이벤트 검색", "이벤트 검색"))
	listCard.QuickReplies.Add(k.QuickReply{}.New("코로나 검색", "코로나 검색"))

	c.JSON(200, listCard.Build())
}
