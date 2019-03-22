package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

const firstJobPageUrl string = "https://www.karriere.at/jobs?page=%d"

func main() {
	fmt.Println("Starting...")

	var jobDetailsChannel chan *JobsDetails = make(chan *JobsDetails)
	go saveResultToFile(jobDetailsChannel)

	start := time.Now()

	noOfPAges := 10
	scrapPage(noOfPAges, jobDetailsChannel)
	//wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}

func scrapPage(noOfPAges int, c chan *JobsDetails) {
	collector := getNewColletcor()

	collector.OnHTML(".m-jobItem__titleLink", func(e *colly.HTMLElement) {
		go collector.Visit(e.Attr("href"))
	})

	collector.OnHTML(".c-jobDetail", func(e *colly.HTMLElement) {
		result := new(JobsDetails)
		result.Url = e.Request.URL.String()
		result.Title = e.DOM.Find(".m-jobHeader__container .m-jobHeader__jobTitle").Text()
		result.Company = e.DOM.Find(".m-jobHeader__container .m-jobHeader__companyName").Text()

		metaItems := e.DOM.Find(".m-jobHeader__container .m-jobHeader__metaList li")
		result.Location = metaItems.First().Text()
		result.Date = metaItems.Last().Text()

		html := e.DOM.Find(".m-jobContent__jobText").Text()
		result.Content = html
		//c <- result
	})

	collector.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	i := 1
	for i <= noOfPAges {
		url := fmt.Sprintf(firstJobPageUrl, i)
		fmt.Println("Scrapping", url)
		collector.Visit(url)
		i = i + 1
	}

	collector.Wait()
}

func saveResultToFile(jobDetailsChannel chan *JobsDetails) {
	os.Remove("result.json")
	f, err := os.OpenFile("result.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	sep := "\n"
	defer f.Close()
	for {
		result := <-jobDetailsChannel
		json, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		//fmt.Println("", result.url, result.company, result.title, result.location, result.date)

		f.WriteString(string(json) + sep)
	}
}

func getNewColletcor() *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
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

type JobsDetails struct {
	Url      string
	Title    string
	Company  string
	Location string
	Date     string
	Content  string
}
