package search

import "time"

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
