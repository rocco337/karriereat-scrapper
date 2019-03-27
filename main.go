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
	//"github.com/arpitgogia/rake"
	"github.com/afjoseph/RAKE.Go"
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

	fullContent:=""
	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

		cleanContent:=strings.ToLower(jobsDetails.Title)
		cleanContent = reg.ReplaceAllString(cleanContent, " ")
		cleanContent = stopwords.CleanString(cleanContent, "de", true)	


		words := strings.Fields(cleanContent)
		sort.Sort(byLength(words))

		stemmed :=""
		for _, word := range words {
			stemmed = stemmed + " " + german.Stemmer.Stem(word)
			fmt.Println(german.Stemmer.Stem(word))
		}
		//fmt.Println(stemmed)
		fullContent = fullContent + " " + stemmed
		line, endOfFile = fileReader.ReadLine()
		//endOfFile = true		
	}
	candidates := rake.RunRake(fullContent)

	for _, candidate := range candidates {
		fmt.Printf("%s --> %f\n", candidate.Key, candidate.Value)
	}
	// rakeResult :=rake.WithText(fullContent)
	// 	for key, value := range rakeResult {
	// 		fmt.Println(key, value)
	// 	}
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
