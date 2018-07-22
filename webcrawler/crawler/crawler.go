package crawler

import (
	"errors"
	"fmt"
	"github.com/TheGUNNER13/playgound-go/webcrawler/analyse"
	"github.com/TheGUNNER13/playgound-go/webcrawler/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type crawl struct {
	filters         []Restriction
	operation       Operation
	waitTime        time.Duration
	visited         map[string]struct{}
	workers         int
	analyser        analyse.Analyze
	visitedLock     sync.RWMutex
	getResponseBody bool
	httpTimeout     time.Duration
}

type Restriction func(*url.URL) bool
type Operation func(*url.URL) *url.URL

type Setting struct {
	Restrictions    []Restriction //these functions can be used to restrict the web crawler, for example, restrict to internal links, respect robots.txt etc.
	Operation       Operation     //this function can take in any url modification that we might want to do
	Workers         int           //number
	WaitTimes       time.Duration
	HttpTimeout     time.Duration
	GetResponseBody bool
}

type Crawler interface {
	Crawl(*url.URL, chan model.CrawlerOutput, chan struct{})
}

func NewCrawler(analyser analyse.Analyze, setting Setting) Crawler {
	if setting.Operation == nil {
		setting.Operation = func(i *url.URL) *url.URL {
			return i
		}
	}
	if analyser == nil {
		analyser = analyse.NewAnalyser(&http.Client{
			Timeout: setting.HttpTimeout,
		})
	}
	if setting.Workers <= 0 {
		setting.Workers = 1
	}
	return &crawl{
		filters:         setting.Restrictions,
		operation:       setting.Operation,
		workers:         setting.Workers,
		analyser:        analyser,
		waitTime:        setting.WaitTimes,
		visited:         make(map[string]struct{}),
		getResponseBody: setting.GetResponseBody,
		httpTimeout:     setting.HttpTimeout,
	}
}

func (c *crawl) Crawl(startURL *url.URL, outWriter chan model.CrawlerOutput, end chan struct{}) {
	defer close(outWriter)

	linkChan := make(chan *url.URL, c.workers*100)
	linkChan <- startURL

	done := make(chan struct{})
	go func() {
		<-end
		log.Println("Received signal to stop")
		for i := 0; i < c.workers; i++ {
			done <- struct{}{}
		}
	}()

	wg := &sync.WaitGroup{}
	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			log.Println("Starting worker #", worker+1)

			for {
				select {
				case <-done:
					log.Println("Stoping worker #", worker+1)
					return
				case link := <-linkChan:
					time.Sleep(c.waitTime)
					pageLinks, outRes, err := c.getPageLinks(link)
					if err != nil {
						continue
					}
					outWriter <- model.CrawlerOutput{URL: link, PageLinks: pageLinks, ResponseBody: outRes}
					for _, link := range pageLinks {
						select {
						case linkChan <- link:
						default:
							log.Println("Channel full, dropping url ", link.String())
						}
					}
				}
			}
		}(i)
	}
	wg.Wait()
	close(linkChan)
	log.Println("All workers stopped now will just process the reamining links on the channel")
	//eat up the links remaining in the channel, not add the links found on the page to the channel
	for link := range linkChan {
		pageLinks, outRes, err := c.getPageLinks(link)
		if err != nil {
			continue
		}
		outWriter <- model.CrawlerOutput{URL: link, PageLinks: pageLinks, ResponseBody: outRes}
	}
	fmt.Println("Stopping Crawler")
}
func (c *crawl) getPageLinks(url *url.URL) ([]*url.URL, io.Reader, error) {
	if c.operation != nil {
		url = c.operation(url)
	}
	if c.runFilters(url) && !c.isVisited(url.String()) {
		return c.analyser.FindLinks(url, c.getResponseBody)
	}
	return nil, nil, errors.New("this url is already visited or is restricted to visit")
}
func (c *crawl) runFilters(url *url.URL) bool {
	for _, filter := range c.filters {
		ok := filter(url)
		if !ok {
			return false
		}
	}
	return true
}

//TODO make the locks performant. Maybe it will be better with a RWMutex?
func (c *crawl) isVisited(url string) bool {
	c.visitedLock.Lock()
	defer c.visitedLock.Unlock()
	_, ok := c.visited[url]
	if !ok {
		c.visited[url] = struct{}{}
	}
	return ok
}
