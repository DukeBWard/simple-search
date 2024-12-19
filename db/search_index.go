package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type SearchIndex struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Value     string
	Urls      []CrawledUrl   `gorm:"many2many:token_urls"`
	CreatedAt *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (s *SearchIndex) TableName() string {
	return "search_index"
}

// maps are cool in go.  map[key]value
// this is saving the index that is being created every hour
func (s *SearchIndex) Save(index map[string][]string, crawledUrls []CrawledUrl) error {
	for value, ids := range index {
		newIndex := &SearchIndex{
			Value: value,
		}
		if err := DBConnect.Where(SearchIndex{Value: value}).FirstOrCreate(newIndex).Error; err != nil {
			return err
		}
		var urlsToAppend []CrawledUrl
		// man i dont like this O(n^2) stuff
		for _, id := range ids {
			for _, url := range crawledUrls {
				if url.ID == id {
					urlsToAppend = append(urlsToAppend, url)
					break
				}
			}
		}
		if err := DBConnect.Model(&newIndex).Association("Urls").Append(&urlsToAppend); err != nil {
			return err
		}
	}
	return nil
}

func (s *SearchIndex) FullTextSearch(value string) ([]CrawledUrl, error) {
	terms := strings.Fields(value)
	var urls []CrawledUrl

	for _, term := range terms {
		var searchIndexes []SearchIndex
		if err := DBConnect.Preload("Urls").Where("value LIKE ?", "%"+term+"%").Find(&searchIndexes).Error; err != nil {
			return nil, err
		}

		for _, searchIndex := range searchIndexes {
			urls = append(urls, searchIndex.Urls...)
		}
	}
	return urls, nil
}
