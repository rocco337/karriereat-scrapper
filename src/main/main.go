package main

import "jobsdataaccess"
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/lib/pq"
)

const firstJobPageUrl string = "https://www.karriere.at/jobs?page=%d"

func main() {
	fmt.Println("Starting...")

	var jobDetailsChannel chan *jobsdataaccess.JobsDetails = make(chan *jobsdataaccess.JobsDetails)
	go saveResultToDb(jobDetailsChannel)

	start := time.Now()

	noOfPAges := 1
	scrapPage(noOfPAges, jobDetailsChannel)

	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}

func scrapPage(noOfPAges int, c chan *jobsdataaccess.JobsDetails) {
	collector := getNewColletcor()

	collector.OnHTML(".m-jobItem__titleLink", func(e *colly.HTMLElement) {
		collector.Visit(e.Attr("href"))
	})

	collector.OnHTML(".c-jobDetail", func(e *colly.HTMLElement) {
		fmt.Println("jobDetail", e.Request.URL.String())
		result := new(jobsdataaccess.JobsDetails)
		result.Url = e.Request.URL.String()
		result.Title = e.DOM.Find(".m-jobHeader__container .m-jobHeader__jobTitle").Text()
		result.Company = e.DOM.Find(".m-jobHeader__container .m-jobHeader__companyName").Text()

		metaItems := e.DOM.Find(".m-jobHeader__container .m-jobHeader__metaList li")
		result.Location = metaItems.First().Text()
		result.Date = metaItems.Last().Text()

		html := e.DOM.Find(".m-jobContent__jobText").Text()
		result.Content = html
		c <- result
	})

	collector.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		panic(err)
	})

	i := 1
	for i <= noOfPAges {
		url := fmt.Sprintf(firstJobPageUrl, i)
		fmt.Println("Scrapping", url)
		collector.Visit(url)
		i = i + 1
		if i%5 == 0 {
			collector.Wait()
		}
	}

	collector.Wait()
}

func saveResultToDb(jobDetailsChannel chan *jobsdataaccess.JobsDetails) {
	jobsDataAccess := new(jobsdataaccess.JobsDataAccess)

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password=postgres dbname=karriereat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	jobsDataAccess.Db = *db
	os.Remove("result.json")
	f, err := os.OpenFile("result.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	sep := "\n"
	for {
		jobDetails := <-jobDetailsChannel

		json, err := json.Marshal(jobDetails)
		if err != nil {
			panic(err)
		}
		go f.WriteString(string(json) + sep)
		go func() {
			jobsDataAccess.SaveJobDetails(jobDetails)
		}()
	}
}
func getNewColletcor() *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
	)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

//https://linuxhint.com/install-pgadmin4-ubuntu/
