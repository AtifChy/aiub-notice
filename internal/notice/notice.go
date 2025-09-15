// Package notice provides functionality to fetch, parse, and cache notices from the AIUB website.
package notice

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "time/tzdata"

	"github.com/AtifChy/aiub-notice/internal/logger"
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
	const maxRetries = 5

	rootURL := "https://www.aiub.edu"
	noticeURL := rootURL + "/category/notices"

	response, err := httpGetWithRetry(noticeURL, maxRetries)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("parse HTML: %w", err)
	}

	document.Find("div.notification").Each(func(_ int, selection *goquery.Selection) {
		title := strings.TrimSpace(selection.Find("h2.title").Text())
		desc := strings.TrimSpace(selection.Find("p.desc").Text())

		dateStr := selection.Find("div.date-custom").Text()
		dateStr = strings.TrimSpace(dateStr)
		re := regexp.MustCompile(`[\s\n]+`)
		dateStr = re.ReplaceAllString(dateStr, " ")

		loc, err := time.LoadLocation("Asia/Dhaka")
		if err != nil {
			logger.L().Error("loading location", slog.String("error", err.Error()))
		}
		date, err := time.ParseInLocation("2 Jan 2006", dateStr, loc)
		if err != nil {
			logger.L().Warn("parsing date",
				slog.String("date", dateStr),
				slog.String("error", err.Error()),
			)
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

	if err := storeCachedNotices(notices); err != nil {
		logger.L().Error("caching notices", slog.String("error", err.Error()))
	}

	return notices, nil
}

func httpGetWithRetry(url string, maxRetries int) (*http.Response, error) {
	var response *http.Response
	var err error

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	_, err = client.Head(url)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	for i := range maxRetries {
		response, err = client.Get(url)
		if err == nil {
			return response, nil
		}

		waitTime := time.Duration(1<<i) * time.Second
		logger.L().Warn(
			"HTTP GET attempt failed",
			slog.Int("attempt", i+1),
			slog.String("error", err.Error()),
			slog.String("wait", waitTime.String()),
		)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("all retries failed: %w", err)
}
