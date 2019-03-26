package main

import (
	"fmt"	
	"karriereat-scrapper/pagescrappers"
	"karriereat-scrapper/dataaccess"
	"karriereat-scrapper/filestorage"

	"time"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting...")

	start := time.Now()

	pageScrapper := new (pagescrappers.PageScrapper)
	pageScrapper.FirstJobPageUrl ="https://www.karriere.at/jobs?page=%d"
	pageScrapper.NoOfPages=5
	pageScrapper.JobDetailsChan = make(chan *dataaccess.JobsDetails)

	//pageScrapper.Init()
	//pageScrapper.ScrapPageRecursively(1)
	fileReader :=new(filestorage.FileReader)
	fileReader.Init("result.json")
	fileReader.ReadAll()

	
	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}




