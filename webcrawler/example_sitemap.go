package main

import (
	"github.com/TheGUNNER13/playgound-go/webcrawler/consumer"
	"github.com/TheGUNNER13/playgound-go/webcrawler/crawler"
	"github.com/TheGUNNER13/playgound-go/webcrawler/model"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var path = "/Users/sdutta/Documents/LBrands/Code/Go/src/github.com/TheGUNNER13/playgound-go/webcrawler/output"

//This example starts crawling "http://www.monzo.com"
//This puts in a filter to not crawl external links like facebook, twitter etc, however the consumer used DOES NOT
// differentiate between internal/external links, so it will list everything given to it by the crawler
//This also puts in an Operation function which strips the urls of '#'
//This functions starts the crawler and lets it crawler for 10 seconds before signalling the crawler to stop.
//This also uses the htmp_sitemap.go consumer to create a simple sitemap
func main() {

	log.SetOutput(os.Stdout)

	toCrawl, _ := url.Parse("http://www.monzo.com")
	var filter crawler.Restriction = func(url *url.URL) bool {
		return url.Host == toCrawl.Host
	}
	var op1 crawler.Operation = func(in *url.URL) *url.URL {
		if in != nil {
			hashIndex := strings.Index(in.String(), "#")
			if hashIndex > 0 {
				out, err := url.Parse(in.String()[:hashIndex])
				if err != nil {
					return in
				}
				return out
			}
		}
		return in
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	out := make(chan model.CrawlerOutput, 100)
	go func() {
		defer wg.Done()
		for each := range out {
			consumer.CreateSiteMap(path, each.URL, each.PageLinks, each.ResponseBody)
		}
	}()
	done := make(chan struct{})

	c := crawler.NewCrawler(nil, crawler.Setting{
		Restrictions:    []crawler.Restriction{filter},
		Operation:       op1,
		WaitTimes:       100 * time.Millisecond,
		Workers:         10,
		GetResponseBody: true,
	})
	go c.Crawl(toCrawl, out, done)

	select {
	case <-time.After(10 * time.Second):
		done <- struct{}{}
	}
	wg.Wait()
}
