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

func scrapPage(pageUrl string, collector *colly.Collector) {
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

	scrapDetailsPage(result.jobLinks)	
}

func scrapDetailsPage(jobLinks []string){
	pageDetailsCollector := getNewColletcor()

	pageDetailsCollector.OnHTML(".c-jobDetail", func(e *colly.HTMLElement) {
		result := new(JobsDetails)	
		result.title = e.DOM.Find(".m-jobHeader__container .m-jobHeader__jobTitle").Text()
		result.company = e.DOM.Find(".m-jobHeader__container .m-jobHeader__companyName").Text()
		
		metaItems:= e.DOM.Find(".m-jobHeader__container .m-jobHeader__metaList li");

		result.location=metaItems.First().Text()
		result.date= metaItems.Last().Text()
		
		html:= e.DOM.Find(".m-jobContent__jobText").Text()
		result.content =html;

		fmt.Println("", result.url, result.company, result.title, result.location, result.date, result.content)
	})

	pageDetailsCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	pageDetailsCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for _, link := range jobLinks {		
		pageDetailsCollector.Visit(link)
		break
	}
}

type JobsPageResult struct {
	jobLinks []string
}

type JobsDetails struct {
	url      string
	title    string
	company  string
	location string
	date     string
	content  string
}
