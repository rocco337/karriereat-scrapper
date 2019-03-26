package pagescrappers

import "karriereat-scrapper/dataaccess"
import "karriereat-scrapper/filestorage"
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/lib/pq"
)

type PageScrapper struct {
	FirstJobPageUrl string
	NoOfPages       int
	JobDetailsChan  chan *dataaccess.JobsDetails
}

func (scrapper *PageScrapper) Init() {
	jobsDataAccess := new(dataaccess.JobsDataAccess)
	jobsDataAccess.Init()

	fileStorage := new(filestorage.FileStorage)
	fileStorage.Init("result.json", true)

	go func() {
		for {
			jobDetails := <-scrapper.JobDetailsChan

			json, err := json.Marshal(jobDetails)
			if err != nil {
				panic(err)
			}

			go jobsDataAccess.SaveJobDetails(jobDetails)
			go fileStorage.AppendLine(string(json))
		}
	}()

}

func (scrapper *PageScrapper) ScrapPageRecursively(currentPage int) {
	collector := scrapper.getCollector()

	collector.OnHTML(".m-jobItem__titleLink", func(e *colly.HTMLElement) {
		collector.Visit(e.Attr("href"))
	})

	collector.OnHTML(".c-jobDetail", func(e *colly.HTMLElement) {
		fmt.Println("jobDetail", e.Request.URL.String())
		result := new(dataaccess.JobsDetails)
		result.Url = e.Request.URL.String()
		result.Title = e.DOM.Find(".m-jobHeader__container .m-jobHeader__jobTitle").Text()
		result.Company = e.DOM.Find(".m-jobHeader__container .m-jobHeader__companyName").Text()

		metaItems := e.DOM.Find(".m-jobHeader__container .m-jobHeader__metaList li")
		result.Location = metaItems.First().Text()
		result.Date = metaItems.Last().Text()

		html := e.DOM.Find(".m-jobContent__jobText").Text()
		result.Content = html
		scrapper.JobDetailsChan <- result
	})

	collector.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		panic(err)
	})

	url := fmt.Sprintf(scrapper.FirstJobPageUrl, currentPage)
	fmt.Println("Scrapping", url)
	collector.Visit(url)

	if currentPage%5 == 0 {
		collector.Wait()
	}

	if currentPage <= scrapper.NoOfPages {
		currentPage = currentPage + 1
		scrapper.ScrapPageRecursively(currentPage)
	}
}

func (scrapper *PageScrapper) getCollector() *colly.Collector {
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
