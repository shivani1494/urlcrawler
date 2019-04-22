package urlcrawler

import (
	"bytes"
	"crypto/tls"
	"github.com/golang/glog"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (u *URLCrawler) getHTMLBodyAndLinks(uri string) []string {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get(uri)
	if err != nil {
		glog.Error(err.Error())
		return []string{}
	}
	defer resp.Body.Close()

	if  !u.mimeTypesSet.Contains(resp.Header.Get("Content-type")) {
		return []string{}
	}

	//fmt.Println("from getDataLinks- " +uri)

	if resp.StatusCode != http.StatusOK {
		glog.Info(uri)
		glog.Error("Failure: ", resp.StatusCode)
		return []string{}
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err.Error())
		return []string{}
	}

	respStr := string(respBytes[:])

	if uri == u.domain {
		u.domainRespStr = respBytes
		return getAllLinks(respStr)
	}

	val := bytes.Compare(respBytes, u.domainRespStr)
	if val == 0 {
		return []string{}
	}

	return getAllLinks(respStr)
}

func getAllLinks(respStr string) []string {

	doc, err := html.Parse(strings.NewReader(respStr))
	if err != nil {
		glog.Error(err.Error())
		return []string{}
	}

	links := make([]string, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					//remove hash and everything after that, refers to
					// the same URL but different parts of the page
					links = append(links, removeHash(a.Val))
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

func parseDomainURL(domain string) []string {

	ret := make([]string, 0)
	u, err := url.Parse(domain)

	if err != nil {
		return ret
	}
	ret = append(ret, u.Scheme)
	ret = append(ret, u.Host)
	ret = append(ret, u.Path)
	return ret
}

func isInternalURL(uri string, domainParts []string) (string, string, bool) {

	if len(domainParts) == 0 {
		return "", "", false
	}

	u, err := url.Parse(uri)
	if err != nil {
		return "", "", false
	}

	//trim the path to get rid of slashes
	path := u.Path

	if u.Path == "" || u.Path == "/" {
		path = ""
	} else if u.Path[0] == '/' {
		path = u.Path[1:]
	}

	if path != "" && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}

	//create an absolute path to return for relative parts
	if u.Host == domainParts[1] || u.Host == "" && u.Path != "" {
		//skip this one because this is the domain URL
		//if the 2 same pages are referred to by diff paths you must ensure
		//to check bytes of the two to verify whether HTML Body is same or not
		if strings.Contains(u.Path, "index") {
			return path, "", false
		}
		uri = domainParts[0] + "://" + domainParts[1] + "/" + path
		return path, uri, true
	}
	return path, uri, false
}

func removeHash(l string) string {
	if strings.Contains(l, "#") {
		var index int
		for n, str := range l {
			if strconv.QuoteRune(str) == "'#'" {
				index = n
				break
			}
		}
		return l[:index]
	}
	return l
}
