package search

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
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
	baseUrl, _ := url.Parse(inputUrl)
	// error or empty
	if err != nil || resp == nil {
		fmt.Println("couldn't crawl this body")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: 0, CrawlData: ParsedBody{}}
	}

	// defer to the end of the call
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("error not 200")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		// resp is html
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			fmt.Println("something went wrong getting data from html body")
			return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
		}
		return CrawlData{Url: inputUrl, Success: true, ResponseCode: resp.StatusCode, CrawlData: data}
	} else {
		// resp is not html
		fmt.Println("error not html")
		return CrawlData{Url: inputUrl, Success: false, ResponseCode: resp.StatusCode, CrawlData: ParsedBody{}}
	}

}

func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return ParsedBody{}, err
	}
	start := time.Now()
	// get links
	links := getLinks(doc, baseUrl)
	// get page title and description
	title, desc := getPageData(doc)
	// get h1 tags
	headings := getHeadings(doc)
	// return time and data
	end := time.Now()
	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Headings:        headings,
		Links:           links,
	}, nil
}

// dfs to get all headers
func getHeadings(node *html.Node) string {
	if node == nil {
		return ""
	}
	var headings strings.Builder
	var findH1 func(*html.Node)
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			if n.FirstChild != nil {
				headings.WriteString(n.FirstChild.Data)
				headings.WriteString(", ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}
	findH1(node)
	// Remove the last comma and space from the concatenated string
	return strings.TrimSuffix(headings.String(), ", ")
}

// returns 2 strings, title and desc
// uses DFS to search through the DOM tree of the html document
func getPageData(node *html.Node) (string, string) {
	if node == nil {
		return "", ""
	}

	title, desc := "", ""
	var findMetaAndTitle func(*html.Node)
	findMetaAndTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild == nil {
				title = ""
			} else {
				title = n.FirstChild.Data
			}
		} else if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}
			if name == "description" {
				desc = content
			}
		}

		// Continue DFS
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findMetaAndTitle(c)
		}
	}

	findMetaAndTitle(node)
	return title, desc
}

// again dfs
func getLinks(node *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if node == nil {
		return links
	}
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					u, err := url.Parse(attr.Val)
					if err != nil ||
						strings.HasPrefix(u.String(), "#") ||
						strings.HasPrefix(u.String(), "mail") ||
						strings.HasPrefix(u.String(), "tel") ||
						strings.HasPrefix(u.String(), "javascript") ||
						strings.HasSuffix(u.String(), ".pdf") ||
						strings.HasSuffix(u.String(), ".md") {
						continue
					}
					if u.IsAbs() {
						// check if same domain
						if isSameHost(u.String(), baseUrl.String()) {
							links.Internal = append(links.Internal, u.String())
						} else {
							links.External = append(links.External, u.String())
						}
					} else {
						// relative link
						rel := baseUrl.ResolveReference(u)
						links.Internal = append(links.Internal, rel.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}
	findLinks(node)
	return links
}

func isSameHost(absoluteUrl string, baseUrl string) bool {
	absUrl, err := url.Parse(absoluteUrl)
	if err != nil {
		return false
	}
	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}
	// checking to see if links are on the same domain
	return absUrl.Host == baseUrlParsed.Host
}
