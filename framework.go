package urlcrawler

import (
	"github.com/deckarep/golang-set"
)

const (
	defaultGoroutines = 2
	defaultDepth      = 2
)

type node struct {
	url   string
	level int
}

type URLCrawler struct {

	//the number of levels to be crawled
	depthN int
	//domainName to crawl
	domain string
	//storing scheme, hostname, path seperately for easier comparision
	domainParts []string
	//keep track of so-far crawled/crawling internal URLs - can convert this into concurrent set
	internalURLSet mapset.Set
	//map[string]bool

	//keep track all external URLs - is concurrent because is a channel
	externalURLs chan string

	doneCrawling  bool
	workerThreads int32
	liveWorkerCt  int32
	domainRespStr []byte

	mimeTypesSet mapset.Set

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


	c.mimeTypesSet = mapset.NewSet()
	c.mimeTypesSet.Add("application/xhtml+xml")
	c.mimeTypesSet.Add("application/xml")
	c.mimeTypesSet.Add( "text/xml")
	c.mimeTypesSet.Add( "text/html")


}
