package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

// AjouLink is the address of where notices of ajou university are being posted
const AjouLink = "https://www.ajou.ac.kr/kr/ajou/notice.do"

// Notice ...
type Notice struct {
	ID       int64  `db:"id" json:"id"`
	Category string `db:"category" json:"category"`
	Title    string `db:"title" json:"title"`
	Date     string `db:"date" json:"date"`
	Link     string `db:"link" json:"link"`
	Writer   string `db:"writer" json:"writer"`
}

// Parse is a function that parses a length of notices
func Parse(url string, length int) []Notice { // doesn't support default value for parameters
	ajouHTML := url
	if url == "" { // As default, use main link
		ajouHTML = fmt.Sprintf("%v?mode=list&articleLimit=%v&article.offset=0", AjouLink, length)
	}

	notices := []Notice{}

	resp, err := soup.Get(ajouHTML)
	if err != nil {
		fmt.Println("[Parser] Check your HTML connection.", err)
		return notices
	}
	doc := soup.HTMLParse(resp)

	ids := doc.FindAll("td", "class", "b-num-box")
	if len(ids) == 0 {
		fmt.Println("[Parser] Check your parser.")
		return notices
	}

	titles := doc.FindAll("div", "class", "b-title-box")
	dates := doc.FindAll("span", "class", "b-date")
	categories := doc.FindAll("span", "class", "b-cate")
	//links := doc.FindAll("div", "class", "b-title-box")
	writers := doc.FindAll("span", "class", "b-writer")
	for i := 0; i < len(ids); i++ {
		id, _ := strconv.ParseInt(strings.TrimSpace(ids[i].Text()), 10, 64)
		title := strings.TrimSpace(titles[i].Find("a").Text())
		link := titles[i].Find("a").Attrs()["href"]
		category := strings.TrimSpace(categories[i].Text())
		date := strings.TrimSpace(dates[i].Text())
		writer := writers[i].Text()

		duplicate := "[" + writer + "]"
		if strings.Contains(title, duplicate) {
			title = strings.TrimSpace(strings.Replace(title, duplicate, "", 1))
		}

		notice := Notice{ID: id, Category: category, Title: title, Date: date, Link: AjouLink + link, Writer: writer}
		notices = append(notices, notice)
	}

	return notices
}

// func CheckConnection() (ok bool) {
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	_, err := http.Get(AjouLink)
// 	if err == nil {
// 		return true
// 	}
// 	return false
// }

/* 어차피 인터넷 연결이 없으면 카톡이 전송 안됨.
if err := models.CheckConnection(); err == false {
	c.JSON(404, models.BuildSimpleText("인터넷 연결 확인 후 잠시 후 다시 시도해보세요."))
	return
}
*/
