package models

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

// AjouLink is the address of where notices of ajou university are being posted
const AjouLink = "https://www.ajou.ac.kr/kr/ajou/notice.do"

// Notice ...
type Notice struct {
	ID     int64  `db:"id" json:"id"`
	Title  string `db:"title" json:"title"`
	Date   string `db:"date" json:"date"`
	Link   string `db:"link" json:"link"`
	Writer string `db:"writer" json:"writer"`
}

// Parse is a function that parses a length of notices
func Parse(url string, length int) []Notice { // doesn't support default value for parameters
	var ajouHTML string = url
	if url == "" { // As default, use main link
		ajouHTML = fmt.Sprintf("%v?mode=list&articleLimit=%v&article.offset=0", AjouLink, length)
	} else {
		ajouHTML = url
	}

	resp, err := soup.Get(ajouHTML)
	if err != nil {
		log.Fatalln("Check your HTML connection.", err)
	}
	doc := soup.HTMLParse(resp)

	notices := []Notice{}

	ids := doc.FindAll("td", "class", "b-num-box")
	if len(ids) == 0 {
		fmt.Println("Check your parser.")
		return notices
	}

	titles := doc.FindAll("div", "class", "b-title-box")
	dates := doc.FindAll("span", "class", "b-date")
	//links := doc.FindAll("div", "class", "b-title-box")
	writers := doc.FindAll("span", "class", "b-writer")
	for i := 0; i < length; i++ {
		id, _ := strconv.ParseInt(strings.TrimSpace(ids[i].Text()), 10, 64)
		title := strings.TrimSpace(titles[i].Find("a").Text())
		link := titles[i].Find("a").Attrs()["href"]
		date := strings.TrimSpace(dates[i].Text())
		writer := writers[i].Text()

		duplicate := "[" + writer + "]"
		if strings.Contains(title, duplicate) {
			writer = strings.TrimSpace(strings.Replace(writer, duplicate, "", 1))
		}

		notice := Notice{ID: id, Title: title, Date: date, Link: AjouLink + link, Writer: writer}
		notices = append(notices, notice)
	}

	return notices
}

func CheckConnection() (ok bool) {
	_, err := http.Get(AjouLink)
	if err != nil {
		return false
	}
	return true
}
