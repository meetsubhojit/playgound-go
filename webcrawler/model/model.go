package model

import (
	"io"
	"net/url"
)

type CrawlerOutput struct {
	URL          *url.URL
	PageLinks    []*url.URL
	ResponseBody io.Reader
}
