package search

import (
	"dukebward/search/db"
	"fmt"
	"time"
)

func RunEngine() {
	fmt.Println("started search engine crawling")
	defer fmt.Println("crawl finished")

	//get settings
	settings := &db.SearchSetting{}
	err := settings.Get()
	if err != nil {
		fmt.Println("something went wrong getting settings")
		return
	}

	// check engine
	if !settings.SearchOn {
		fmt.Println("search is turned off")
		return
	}
	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCrawlUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("something is not working with the crawling")
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdatedUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Headings:        result.CrawlData.Headings,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Println("something went wrong when updating url")
			}
			continue
		}
		err := next.UpdatedUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageTitle,
			PageDescription: result.CrawlData.PageDescription,
			Headings:        result.CrawlData.Headings,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Println("something went wrong when updating url")
		}
		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // range end
	// if we dont want to add new urls to system
	if !settings.AddNew {
		return
	}
	// if not, add new urls
	for _, newUrl := range newUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Println("something went wrong with adding new url to db")
		}
	}
	fmt.Printf("\n Added %d new urls to the db", len(newUrls))
}

func RunIndex() {
	fmt.Println("search index start")
	defer fmt.Println("search index stopped")

	// get index settings from db
	crawled := &db.CrawledUrl{}
	// get non indexed urls
	notIndexed, err := crawled.GetNotIndex()
	if err != nil {
		return
	}

	// make new index and not non indexed urls
	idx := make(Index)
	idx.Add(notIndexed)
	searchIndex := &db.SearchIndex{}
	// save index to db
	err = searchIndex.Save(idx, notIndexed)
	if err != nil {
		return
	}
	// save urls indexed to true after we index them
	err = crawled.SetIndexedTrue(notIndexed)
	if err != nil {
		return
	}

}
