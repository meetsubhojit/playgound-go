package crawler

import (
	"github.com/TheGUNNER13/playgound-go/webcrawler/analyse"
	"io"
	"net/url"
	"testing"
)

func TestCrawl_getPageLinks(t *testing.T) {
	mockLinks := []*url.URL{urlMaker("a")}
	c := NewCrawler(&mockAnalyse{outlinks: mockLinks}, Setting{})
	crawl, ok := c.(*crawl)
	if !ok {
		t.Fail()
	}
	actualPageLink, _, _ := crawl.getPageLinks(urlMaker("a"))
	if len(actualPageLink) != len(mockLinks) {
		t.Fail()
	}
}

type mockAnalyse struct {
	analyse.Analyze
	outlinks []*url.URL
	outError error
}

func (m *mockAnalyse) FindLinks(*url.URL, bool) ([]*url.URL, io.Reader, error) {
	if m.outError != nil {
		return nil, nil, m.outError
	}
	return m.outlinks, nil, nil
}
func Test_isVisited(t *testing.T) {
	c := crawl{
		visited: map[string]struct{}{"a": {}, "b": {}},
	}
	if !c.isVisited("a") {
		t.Fatalf("Expected true got false")
	}
	if !c.isVisited("b") {
		t.Fatalf("Expected true got false")
	}
	if c.isVisited("c") {
		t.Fatalf("Expected false got true")
	}
	if !c.isVisited("c") {
		t.Fatalf("Expected true got false")
	}
}
func Test_runFilters(t *testing.T) {
	var tests = []struct {
		name         string
		restrictions []Restriction
		out          bool
	}{
		{
			name: "no restriction",
			out:  true,
		},
		{
			name: "one restriction, give true",
			restrictions: []Restriction{func(url *url.URL) bool {
				return true
			}},
			out: true,
		},
		{
			name: "two restriction, give true, false resp.",
			restrictions: []Restriction{func(url *url.URL) bool {
				return true
			}, func(url *url.URL) bool {
				return false
			}},
			out: false,
		},
	}
	c := &crawl{}
	for _, test := range tests {
		c.filters = test.restrictions
		actual := c.runFilters(nil)
		if actual != test.out {
			t.Fatalf("Test %s failed because expecting output to be %v but got %v", test.name, test.out, actual)
		}
	}
}
func urlMaker(link string) *url.URL {
	u, _ := url.Parse(link)
	return u
}
