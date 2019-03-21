package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
)

const firstJobPageUrl string = "https://www.karriere.at/jobs?page=%d"

func main() {

	scrapPage(fmt.Sprintf(firstJobPageUrl, 1), getNewColletcor())

}

func getNewColletcor() *colly.Collector {
	c := colly.NewCollector()
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	return c
}

func scrapPage(pageUrl string, collector *colly.Collector) *JobsPageResult {
	result := new(JobsPageResult)

	collector.OnHTML(".m-jobItem__titleLink", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		result.jobLinks = append(result.jobLinks, link)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("\nVisiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.Visit(pageUrl)

	for _, link := range result.jobLinks {
		pageDetailsCollector := getNewColletcor()
		scrapJobDetailsPage(link, pageDetailsCollector)
	}

	return result
}
func scrapJobDetailsPage(url string, collector *colly.Collector) *JobsDetails {
	result := new(JobsDetails)
	result.url = url

	collector.OnHTML(".m-jobHeader__jobTitle", func(e *colly.HTMLElement) {
		result.title = e.Text
	})

	collector.OnHTML(".m-jobHeader__companyName", func(e *colly.HTMLElement) {
		result.company = e.Text
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("\nVisiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.Visit(url)

	fmt.Println("", result.url, result.company, result.title)
	return result
}

type JobsPageResult struct {
	jobLinks []string
}

type JobsDetails struct {
	url      string
	title    string
	company  string
	location string
	date     time.Time
	content  string
}
