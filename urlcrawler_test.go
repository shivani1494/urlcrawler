package urlcrawler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTest(url string) URLCrawler {
	var urlC URLCrawler
	urlC.NewURLCrawler(url)
	return urlC
}

//to test out all the other cases just set it equal to u1
var (
	u1  = "http://www.chrislengerich.com/"
	u2  = "http://eepurl.com/cxPQd5"
	u3  = "http://homestarrunner.com/"
	u4  = "http://www.wayback.com/"
	u5  = "https://www.badgerbadgerbadger.com/"
	u6  = "http://www.taco.com/"
	u7  = "http://www.impactsearch.io/"
	u8  = "http://emerginstars.com/"
	u9  = "http://www.samvitjain.com/./assets"
	u10 = "http://...com"
)

type interalURLTest struct {
	input    string
	internal string
}

var isInterURLTests = []interalURLTest{

	{"http://www.chrislengerich.com/", "essay.html"},
	{"http://www.chrislengerich.com/", "verified.html"},
	{"http://www.chrislengerich.com/", "http://www.zhenfund.com/"},
	{"http://eepurl.com/cxPQd5", "https://medium.com"},
}

func Test_CrawlDomainURL(t *testing.T) {

	urlC := setupTest(u1)

	fmt.Println("======= Crawling for: " + u1 + " =======")
	err := urlC.CrawlDomainURL()

	if err != nil {
		assert.Error(t, err, "could not crawl the domain- "+u1)
	}

	fmt.Println("======= Printing Status for: " + u1 + " =======")
	urlC.GetStatus()

	fmt.Println("======= Printing Result for: " + u1 + " =======")
	urlC.GetResult()

}

func Test_getHTMLBodyAndLinks(t *testing.T) {

	urlC := setupTest(u1)
	links := urlC.getHTMLBodyAndLinks(u1)
	for i := range links {
		fmt.Println(links[i])
	}
}

func Test_parseDomainURL(t *testing.T) {

	links := parseDomainURL(u1)
	fmt.Println("======= Parsed Domain URL for: " + u1 + " =======")
	for i := range links {
		fmt.Println(links[i])
	}
}

func Test_isInternalURL(t *testing.T) {

	for _, tt := range isInterURLTests {
		_, _, isInternal := isInternalURL(tt.internal, parseDomainURL(tt.input))
		fmt.Print("======= base URL-- " + tt.input + " secondary URL-- " +
			tt.internal + " isInternal? ")
		fmt.Print(isInternal)
		fmt.Println()
	}
}
