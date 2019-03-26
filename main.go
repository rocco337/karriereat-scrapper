package main

import (
	"encoding/json"
	"fmt"
	"karriereat-scrapper/dataaccess"
	"karriereat-scrapper/filestorage"
	"karriereat-scrapper/pagescrappers"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bbalet/stopwords"
	"github.com/dchest/stemmer/german"
)

func main() {
	fmt.Println("Starting...")

	start := time.Now()

	pageScrapper := new(pagescrappers.PageScrapper)
	pageScrapper.FirstJobPageUrl = "https://www.karriere.at/jobs?page=%d"
	pageScrapper.NoOfPages = 5
	pageScrapper.JobDetailsChan = make(chan *dataaccess.JobsDetails)

	//pageScrapper.Init()
	//pageScrapper.ScrapPageRecursively(1)

	//languagedetector := new(languagedetector.LanguageDetector)
	//languagedetector.Init()
	//languagedetector.DetectAndSave()

	fileReader := new(filestorage.FileReader)
	fileReader.Init("result.json")

	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

		cleanContent := stopwords.CleanString(jobsDetails.Content, "de", true)
		words := strings.Fields(reg.ReplaceAllString(cleanContent, " "))
		sort.Sort(byLength(words))

		for _, word := range words {
			fmt.Println(german.Stemmer.Stem(word))
		}

		//line, endOfFile = fileReader.ReadLine()
		endOfFile = true
	}
	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}

type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
