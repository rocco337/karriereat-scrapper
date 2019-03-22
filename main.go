package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	var jobDetailsChannel chan *JobsDetails = make(chan *JobsDetails)
	go saveResultToDb(jobDetailsChannel)

	start := time.Now()

	noOfPAges := 500
	scrapPage(noOfPAges, jobDetailsChannel)

	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}

func scrapPage(noOfPAges int, c chan *JobsDetails) {
	collector := getNewColletcor()

	collector.OnHTML(".m-jobItem__titleLink", func(e *colly.HTMLElement) {
		collector.Visit(e.Attr("href"))
	})

	collector.OnHTML(".c-jobDetail", func(e *colly.HTMLElement) {
		fmt.Println("jobDetail", e.Request.URL.String())
		result := new(JobsDetails)
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

func saveResultToDb(jobDetailsChannel chan *JobsDetails) {

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres "+
		"password=postgres dbname=karriereat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	os.Remove("result.json")
	f, err := os.OpenFile("result.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	sep := "\n"
	for {
		result := <-jobDetailsChannel

		json, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		go f.WriteString(string(json) + sep)
		go func() {
			_, err := db.Exec(`INSERT INTO jobs(url, title, company,location, date, content)
	VALUES($1,$2,$3,$4,$5,$6)`, result.Url, result.Title, result.Company, result.Location, result.Date, result.Content)
			if err != nil {
				panic(err)
			}
		}()
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

//https://linuxhint.com/install-pgadmin4-ubuntu/
