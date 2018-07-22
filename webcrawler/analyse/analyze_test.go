package analyse

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAnalyze_FindLinks(t *testing.T) {

	var respHTML string
	var respErr error
	makeGetCall = func(analyzer *analyze, url string) (*http.Response, error) {
		if respErr != nil {
			return nil, respErr
		}
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBuffer([]byte(respHTML))),
		}, nil
	}

	var tests = []struct {
		name         string
		respHTML     string
		respErr      error
		inLink       string
		needResponse bool
		outLink      []string
		outErr       error
	}{
		{
			name:     "one url is found",
			inLink:   "www.test1.com",
			respHTML: `<a href="http://www.test1.com">Test1</a>`,
			outLink:  []string{"http://www.test1.com"},
		},
		{
			name:     "one url is found, correction needed",
			inLink:   "http://www.test1.com",
			respHTML: `<a href="./a/test.html">Test1</a>`,
			outLink:  []string{"http://www.test1.com/a/test.html"},
		},
		{
			name:    "http get returns error",
			inLink:  "www.test1.com",
			respErr: errors.New("some error"),
			outErr:  errors.New("some error"),
		},
	}
	analyser := NewAnalyser(nil)
	for _, test := range tests {
		respHTML = test.respHTML
		respErr = test.respErr
		actualLinks, _, actualErr := analyser.FindLinks(urlMaker(test.inLink), test.needResponse)
		if test.outErr != nil {
			if actualErr == nil || !strings.Contains(test.outErr.Error(), actualErr.Error()) {
				t.Fatalf("TEST %s failed becuase expected error : %v, did not match actual error : %v", test.name, test.respErr, actualErr)
			}
		} else {
			if len(test.outLink) != len(actualLinks) {
				t.Fatalf("TEST %s failed becuase length of expected links : %d, did not match actual length : %d", test.name, len(test.outLink), len(actualLinks))
			} else if !matchLinks(actualLinks, test.outLink) {
				t.Fatalf("TEST %s failed becuase expected links : %v, did not match actual links : %v", test.name, test.outLink, actualLinks)
			}
		}
	}
}
func Test_correctTheLink(t *testing.T) {
	var tests = []struct {
		name   string
		in     string
		domain string
		out    string
	}{
		{
			name:   "Test1",
			in:     "hello/",
			domain: "http://localhost:5100",
			out:    "http://localhost:5100/hello/",
		},
		{
			name:   "Test2",
			in:     "./../d",
			domain: "http://localhost:5100/a/b/c/",
			out:    "http://localhost:5100/a/b/d",
		},
		{
			name:   "Test3",
			in:     "//localhost:5100/a",
			domain: "http://localhost:5100/a/b/c/",
			out:    "http://localhost:5100/a",
		},
		{
			name:   "Test4",
			in:     "./d",
			domain: "http://localhost:5100/a/b/c/",
			out:    "http://localhost:5100/a/b/c/d",
		},
	}
	for _, test := range tests {
		actual := correctTheLink(test.in, urlMaker(test.domain)).String()
		if test.out != actual {
			t.Fatalf("Test %s failed because expected url %s does not match actual url %s", test.name, test.out, actual)
		}
	}
}
func matchLinks(urls []*url.URL, strLinks []string) bool {
	for i, s := range strLinks {
		if s != urls[i].String() {
			return false
		}
	}
	return true
}
func urlMaker(link string) *url.URL {
	u, _ := url.Parse(link)
	return u
}
