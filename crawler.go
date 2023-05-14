package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
	"time"
)

type HugoThemeCrawler struct {
	c                *colly.Collector
	crawlErrorChan   chan error // When an error occurs, the error will be put in this channel.
	finish           chan byte  // When the crawler ends gracefully the end byte will be put in this channel.
	headers          map[string]string
	maxTry           int
	crawledUrlCount  int
	crawledItemCount int
	themeList        []map[string]interface{}
}

func GetCrawler() *HugoThemeCrawler {
	return &HugoThemeCrawler{}
}

func (ht *HugoThemeCrawler) InitHugoThemeCrawler() {
	ht.headers = make(map[string]string)
	ht.headers["user-agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	ht.headers["accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	ht.headers["authority"] = "themes.gohugo.io"
	ht.c = colly.NewCollector(
	//colly.AllowURLRevisit(),
	)

	// Set timeout 5s
	ht.c.SetRequestTimeout(5 * time.Second)

	// When the crawler is finished it will call OnScraped()
	ht.c.OnScraped(func(r *colly.Response) {
		if ht.crawledUrlCount == ht.crawledItemCount {
			log.Printf("In startCrawlHugoThemes: The crawler ends gracefully. Crawled %d urls! Crawled %d items!\n", ht.crawledUrlCount, ht.crawledItemCount)
			ht.finish <- 1
		}
	})

	// When a crawler error occurs it will call OnError()
	ht.c.OnError(func(r *colly.Response, err error) {
		log.Printf("In OnError: The crawler ends with an error. -> %s\n", err)
		ht.crawlErrorChan <- err
	})

	// When a matching item is found it will call OnHTML()
	ht.c.OnHTML(".ma0.sans-serif.bg-primary-color-light", func(e *colly.HTMLElement) {
		if ht.parserSelector(e.Request.URL.String()) {
			urls := ht.parseTagsPage(e.DOM)
			ht.crawledUrlCount += len(urls)
		} else {
			details := ht.parseDetailsPage(e.DOM, e.Request.URL.String())
			if len(details) != 0 {
				ht.themeList = append(ht.themeList, details)
				ht.crawledItemCount += 1
			}
		}
	})

	// When a request is sent it will call OnRequest()
	ht.c.OnRequest(func(r *colly.Request) {
		for key, value := range ht.headers {
			r.Headers.Set(key, value)
		}
	})

	ht.maxTry = 5
}

func (ht *HugoThemeCrawler) startCrawlHugoThemes() error {
	var crawlError error

	// Before crawling task starts, initialize these channels
	ht.crawlErrorChan = make(chan error)
	ht.finish = make(chan byte)

	var startUrl = "https://themes.gohugo.io"
	go ht.c.Visit(startUrl)
	log.Printf("In startCrawlHugoThemes: get url -> %s", startUrl)

	for {
		var breakFlag = false
		select {
		case crawlError = <-ht.crawlErrorChan:
			breakFlag = true
		case <-ht.finish:
			breakFlag = true
		default:
		}
		if breakFlag {
			break
		}
	}
	return crawlError
}

func (ht *HugoThemeCrawler) parseDetailsPage(DOM *goquery.Selection, orginalUrl string) map[string]interface{} {
	details := make(map[string]interface{})
	var value any
	details["url"] = orginalUrl
	DOM.Find("li.mb2").Each(func(_ int, li *goquery.Selection) {
		label := li.Find("span.label").Text()
		if label == "Author:" {
			value = li.Find("a").Text()
		} else {
			value = li.Find("span.value").Text()
		}
		details[label] = value
	})
	var tags []string
	DOM.Find(".mb2.mt4 a").Each(func(_ int, a *goquery.Selection) {
		text := strings.TrimSpace(a.Text())
		tags = append(tags, text)
	})
	details["Tags:"] = tags
	//log.Println("In parseDetailsPage: Parsing item -> ", details)
	return details
}

func (ht *HugoThemeCrawler) parseTagsPage(DOM *goquery.Selection) []string {
	var hrefs []string
	DOM.Find(".link.db.shadow-hover.gray.mb4.w-100.w-30-ns").Each(func(_ int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {
			hrefs = append(hrefs, href)
			go ht.c.Visit(href)
			//log.Printf("In parseTagsPage: Parsing href -> %s\n", href)
		}
	})
	return hrefs
}

func (ht *HugoThemeCrawler) parserSelector(link string) bool {
	return !strings.Contains(link, "themes/")
}

func (ht *HugoThemeCrawler) getThemes() []map[string]interface{} {
	return ht.themeList
}
