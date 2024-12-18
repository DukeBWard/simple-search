package search

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CrawlData struct {
	Url          string `json:"url"`
	Success      bool   `json:"success"`
	ResponseCode int    `json:"responseCode"`
	CrawlData    ParsedBody
}

type ParsedBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Headings        string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

func runCrawl(inputUrl string) CrawlData {
	resp, err := http.Get(inputUrl)
	baseUrl, err := url.Parse(inputUrl)
	// error or empty
	if err != nil || resp == nil {
		fmt.Println("couldnt crawl this body")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: 0, CrawlData: ParsedBody{}}
	}

	// defer to the end of the call
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("error not 200")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}
	content := resp.Header.Get("Content-Type")
	if strings.HasPrefix(co)

}
