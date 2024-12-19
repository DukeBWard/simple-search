package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Url             string         `json:"url" gorm:"unique;not null"`
	Success         bool           `json:"success" gorm:"default:null"`
	CrawlDuration   time.Duration  `json:"crawlDuration"`
	ResponseCode    int            `json:"responseCode" gorm:"type:smallint"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Headings        string         `json:"headings"`
	LastTested      *time.Time     `json:"lastTested"` // Use pointer so this value can be nil
	Indexed         bool           `json:"indexed" gorm:"default:false"`
	CreatedAt       *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// saving crawled urls to the database
func (crawled *CrawledUrl) UpdatedUrl(input CrawledUrl) error {
	tx := DBConnect.Select("url", "success", "crawl_duration", "response_code", "page_title", "page_description",
		"headings", "last_tested", "indexed", "updated_at").Omit("created_at").Save(&input)
	if tx.Error != nil {
		fmt.Print(tx.Error)
		return tx.Error
	}
	return nil
}

// receiver, parameters, return values
func (crawled *CrawledUrl) GetNextCrawlUrls(limit int) ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConnect.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

func (crawled *CrawledUrl) Save() error {
	tx := DBConnect.Save(&crawled)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (crawled *CrawledUrl) GetNotIndex() ([]CrawledUrl, error) {
	var urls []CrawledUrl
	// tx is transaction
	tx := DBConnect.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

func (crawled *CrawledUrl) SetIndexedTrue(urls []CrawledUrl) error {
	for _, url := range urls {
		url.Indexed = true
		tx := DBConnect.Save(&url)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
