package analyse

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	aTag    = "a"
	hrefTag = "href"
)

var makeGetCall = func(analyser *analyze, url string) (*http.Response, error) {
	return analyser.httpclient.Get(url)
}

//Analyze provides an interface, which given a url link will
//1- Get the html page
//2- Find for links in the html page
//3- Correct the link, i.e. correct if according to the domain if its a relative path
//4- Return the found links, along with the http responseBody is required.
type Analyze interface {
	FindLinks(*url.URL, bool) ([]*url.URL, io.Reader, error)
}
type analyze struct {
	httpclient *http.Client
}

func NewAnalyser(client *http.Client) Analyze {
	if client == nil {
		client = http.DefaultClient
	}
	return &analyze{
		httpclient: client,
	}
}
func (a *analyze) FindLinks(in *url.URL, needResBody bool) ([]*url.URL, io.Reader, error) {
	var out []*url.URL

	res, err := makeGetCall(a, in.String())
	if err != nil {
		log.Printf("HTML error at in %s, error : %s", in, err.Error())
		return nil, nil, err
	}
	if res == nil {
		return nil, nil, errors.New("Got empty response for url : " + in.String())
	}
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	var copy1, copy2 *bytes.Buffer
	copy1 = bytes.NewBuffer(b)
	if needResBody {
		copy2 = bytes.NewBuffer(b)
	}
	htmlBody := html.NewTokenizer(copy1)
	for {
		tokenType := htmlBody.Next()
		switch tokenType {
		case html.ErrorToken:
			if htmlBody.Err() != io.EOF {
				log.Printf("HTML error at in %s, error : %s", in, htmlBody.Err().Error())
			}
			return out, copy2, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := htmlBody.Token()
			if token.DataAtom.String() == aTag {
				for _, attr := range token.Attr {
					if attr.Key == hrefTag {
						link := correctTheLink(attr.Val, in)
						if link != nil {
							out = append(out, link)
						}
					}
				}
			}
		}
	}
}

func correctTheLink(in string, domain *url.URL) *url.URL {
	inUri, err := url.Parse(in)
	if err != nil {
		return nil
	}
	return domain.ResolveReference(inUri)
}
