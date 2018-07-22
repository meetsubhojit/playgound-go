package consumer

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

var htmlBodyPattern = "<html><header><title>%s</title></header><body>%s</body></html>"
var list = "<ul><lh>Links on this page</lh>%s</ul>"

var alreadyProcessed = map[string]struct{}{}

//This function creates a simple sitemaps using the links, response bodies given by the crawler.
//It sees if the page is already rendered, if not, it creates an html page having the links present on this page in a unordered list
//TODO make it process parallely is required, this function for now processes each page serially, so its not optimized
//TODO make the HTML page prettier
func CreateSiteMap(path string, url *url.URL, links []*url.URL, resBody io.Reader) {
	if _, ok := alreadyProcessed[url.String()]; !ok {
		alreadyProcessed[url.String()] = struct{}{}

		log.Printf("Rendering URL %s\n", url.String())
		filepath := getFullFileName(path, url)
		title := getHeaderTitle(resBody)

		outLinks := ""
		for _, link := range links {
			filepath := getFullFileName(path, link)
			outLinks = outLinks + fmt.Sprintf("<li><a href=%s>%s</a>", filepath, link.String())
		}

		outHtml := fmt.Sprintf(htmlBodyPattern, title, fmt.Sprintf(list, outLinks))
		err := ioutil.WriteFile(filepath, []byte(outHtml), 0644)
		if err != nil {
			log.Println("Error rendering url ", url.String(), " error : ", err.Error())
		}
	}
}
func getFullFileName(path string, url *url.URL) string {
	name := hash(url.String())
	return fmt.Sprintf("%s/%s.html", path, name)
}

//TODO optimise this function to give a consistent output for a particular URL. Maybe its doing right now?
func hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

//TODO optimize this, string extracting maybe the fastest?
func getHeaderTitle(resBody io.Reader) string {
	b, err := ioutil.ReadAll(resBody)
	if err != nil {
		return ""
	}
	stringB := string(b)
	start := strings.Index(stringB, `<title>`)
	end := strings.Index(stringB, `</title>`)
	if start < 0 || end < 0 {
		start = strings.Index(stringB, `<TITLE>`)
		end = strings.Index(stringB, `</TITLE>`)
	}
	if start < 0 || end < 0 {
		return ""
	}
	return stringB[start+7 : end]
}
