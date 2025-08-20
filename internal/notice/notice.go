// Package notice provides functionality to fetch, parse, and cache notices from the AIUB website.
package notice

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Notice struct {
	Date  time.Time
	Title string
	Desc  string
	Link  string
}

func GetNotices() ([]Notice, error) {
	var notices []Notice

	rootURL := "https://www.aiub.edu"
	noticeURL := rootURL + "/category/notices"

	response, err := http.Get(noticeURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	document.Find("div.notification").Each(func(_ int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find("h2.title").Text())
		desc := strings.TrimSpace(selection.Find("p.desc").Text())

		dateStr := selection.Find("div.date-custom").Text()
		dateStr = strings.TrimSpace(dateStr)
		re := regexp.MustCompile(`[\s\n]+`)
		dateStr = re.ReplaceAllString(dateStr, " ")

		loc, _ := time.LoadLocation("Asia/Dhaka")
		date, err := time.ParseInLocation("2 Jan 2006", dateStr, loc)
		if err != nil {
			log.Printf("Error: parsing date '%s': %v", dateStr, err)
		}

		link, _ := selection.Find("a").Attr("href")
		link = rootURL + link

		notices = append(notices, Notice{
			Date:  date,
			Title: title,
			Desc:  desc,
			Link:  link,
		})
	})

	if err := saveNoticesCache(notices); err != nil {
		log.Printf("Error: failed to cache notices: %v", err)
	}

	return notices, nil
}
