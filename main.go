package main 
import "fmt"
import "time"
import "github.com/gocolly/colly"

const firstJobPageUrl string ="https://www.karriere.at/jobs?page=1"

func main() {	
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
		}

	pageResult := scrapPage(firstJobPageUrl, 1)
	for index, link := range pageResult.jobLinks {
		fmt.Printf("# %d. %v", index, link)	
	}
}

func scrapPage(pageUrl string, pageNumber int) *JobsPageResult {
	fmt.Printf("Scrapping page # %d", pageNumber)
	result :=new(JobsPageResult)
	return result
}
func scrapJobDetailsPage(url string) *JobsDetails {
	result :=new(JobsDetails)
	return result
}



type JobsPageResult struct {
	jobLinks []string
}

type JobsDetails struct{
	url string
	title string
	company string
	location string
	date time.Time
	content string
}