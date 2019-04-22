package urlcrawler

import (

	"sync/atomic"
	"fmt"
	"net/url"
	"time"
	"errors"
	"github.com/golang/glog"
)

func (u *URLCrawler) GetStatus() {

	for {
		if u.doneCrawling {
			break
		} else {
			time.Sleep(700 * time.Millisecond)
			fmt.Println("Completing- ", u.internalURLSet.Cardinality())
		}
	}
	fmt.Println("Completed- ", u.internalURLSet.Cardinality())
}

func (u *URLCrawler) GetResult() {

	result := make(map[string]int)

	if len(u.externalURLs) == 0 {
		fmt.Println("No External URLs found")
	}

	for currURL := range u.externalURLs {
		parsedURL, _ := url.Parse(currURL)
		if parsedURL.Host == "" {
			continue
		}
		if _, ok := result[parsedURL.Host]; !ok {
			result[parsedURL.Host] = 1
			continue
		}
		result[parsedURL.Host]++
	}

	for i := range result {
		fmt.Println(i, "-->" ,result[i])
	}
}

func (u *URLCrawler) crawlCurrentURL(queueNodes chan *node) {

	for {

		if len(queueNodes) == 0 {
			break
		}

		i := <- queueNodes
		path, fmtURL, ok := isInternalURL(i.url, u.domainParts)

		//if external link
		if !ok {
			glog.Info(i.url)
			u.externalURLs <- i.url
			continue
		}

		//if page already visited
		//u.mux.Lock()
		if ok := u.internalURLSet.Contains(path); ok {
			glog.Info("from crawlCurrentURL -- visited -- " + fmtURL + " " + path)
			//u.mux.Unlock()
			continue
		}
		//u.mux.Unlock()

		//u.internalURLSet[path] = true
		glog.Info("from crawlCurrentURL -- unvisited -- " + fmtURL + " " + path)
		added := u.internalURLSet.Add(fmtURL)

		if added {

			links := make([]string, 0)

			links = u.getHTMLBodyAndLinks(fmtURL)

			glog.Info("crawled url - " + fmtURL)

			for j := range links {
				if i.level < u.depthN {
					queueNodes <- &node{url: links[j], level: i.level + 1}
				}
			}
		}

		var num int32
		if int32(len(queueNodes)) > 0 || (u.workerThreads - atomic.LoadInt32(&u.liveWorkerCt) > 0) {
			for num = 0; num < (u.workerThreads - atomic.LoadInt32(&u.liveWorkerCt)); num++ {
				atomic.AddInt32(&u.liveWorkerCt, int32(1))
				go u.crawlCurrentURL(queueNodes)
			}
		}
	}

	atomic.AddInt32(&u.liveWorkerCt, int32(-1))
}

func (u *URLCrawler) CrawlDomainURL() error {

	if len(u.domainParts) == 0 {
		return errors.New("could not parse the url " + u.domain)
	}

	queueNodes := make(chan *node, 1000) //initial cap?
	queueNodes <- &node{url:u.domain, level:0}

	atomic.AddInt32(&u.liveWorkerCt, int32(1))

	go u.crawlCurrentURL(queueNodes)
	go u.GetStatus()

	for {
		if atomic.LoadInt32(&u.liveWorkerCt) == int32(0) {
			// this can be made atomic,
			// because I am reading this value in a diff thread
			u.doneCrawling = true
			close(queueNodes)
			close(u.externalURLs)
			break
		}
	}
	//this is like a wait when one thread returns back you can't just let the
	//main thread to end, you have to keep the process going

	return nil
}
