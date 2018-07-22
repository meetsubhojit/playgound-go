# webcrawler

This web crawler has mainly two part, Crawler and Analyser.

The core part if the Crawler (package). The crawler package starts with a given url and it gives the url to an "Analyser" interface which inturn gives it back links on the current url.
The Crawler has many options, like http timeouts, how much time to wait before each crawl etc.
It takes in a list of Restrictions to see if the current URL needs to be crawled. Here ideally we can put restrictions like respecting robots.txt, only internal links etc. 
The crawler accepts an Analyser, so ideally the consumer of the library can write their own Analysers.

The Analyser here is a default one used in the crawler package.
Its main functions is to parse the html tags and extract the links. Additionally is can give back the user the whole HTML page if some post processing is required.

Additionally there is a consumer package, which ideally should be outside the webcrawler, but it is kept for demonstration purposes.
The consumer package listens to the channel where the crawler output its findings, and the consumer can now index this url.
Ideally here we can write anything form making simple sitemaps, storing the data in a DB, write it to files etc.

just "go run example_sitemap.go" to get a simple sitemap for monzo.com

