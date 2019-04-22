package urlcrawler

import (
	"github.com/deckarep/golang-set"
)

const (
	defaultGoroutines = 1
	defaultDepth = 2
)

type node struct {
	url string
	level int
}

type URLCrawler struct {

	//the number of levels to be crawled
	depthN int

	//domainName to crawl
	domain string

	domainParts []string

	//keep track of so-far crawled/crawling internal URLs - can convert this into concurrent set
	internalURLSet mapset.Set
	//map[string]bool

	//keep track all external URLs - is concurrent because is a channel
	externalURLs chan string

	doneCrawling bool
	workerThreads int32
	liveWorkerCt int32
	domainRespStr []byte

	//mux sync.Mutex

}

func (c *URLCrawler) NewURLCrawler(domain string) {

	c.domain = domain
	c.domainParts = parseDomainURL(c.domain)

	c.liveWorkerCt = 0
	c.doneCrawling = false

	c.internalURLSet = mapset.NewSet()

	c.externalURLs = make(chan string, 100)

	c.workerThreads = defaultGoroutines
	c.depthN = defaultDepth
}