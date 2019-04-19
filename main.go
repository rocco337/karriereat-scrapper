package main
import "karriereat-scrapper/filestorage"
import (
	"encoding/json"
	"fmt"
	"karriereat-scrapper/dataaccess"
	"karriereat-scrapper/languagedetector"
	"karriereat-scrapper/pagescrappers"
	"karriereat-scrapper/phraseextractor"
	"time"
    "github.com/dchest/stemmer/german"
	rake "github.com/afjoseph/RAKE.Go"
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

	languagedetector := new(languagedetector.LanguageDetector)
	languagedetector.Init()
	//languagedetector.DetectAndSave()

	fileReader := new(filestorage.FileReader)
	fileReader.Init("result.json")

	phraseextractor:=new (phraseextractor.Phraseextractor)

	// fileStorage:= new(filestorage.FileStorage)
	// fileStorage.Init("phrases.json", true)

	// mappedPhrases := phraseextractor.ExtractPhrasesFromFile(fileReader, languagedetector)
	// phraseextractor.WritePhrasesToFile(mappedPhrases, fileStorage)

	
	phrasesFileReader := new(filestorage.FileReader)
	phrasesFileReader.Init("phrases.json")
	
	mappedPhrases := phraseextractor.ReadPhrasesFromFile(phrasesFileReader)	
	extractKeywordsFromFiles(fileReader, languagedetector, mappedPhrases)
	
	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}


func extractKeywordsFromFiles(fileReader *filestorage.FileReader,languagedetector *languagedetector.LanguageDetector,mappedPhrases map[string]int){
	i:=0
	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		line, endOfFile = fileReader.ReadLine()
	
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		fmt.Println("----",jobsDetails.Url, jobsDetails.Title)
		lang := languagedetector.Detect(jobsDetails.Content)
		if lang == "de" {
			ger := german.Stemmer			
			stemmed:= ger.Stem(jobsDetails.Content)
			candidates := rake.RunRakeI18N(stemmed, phraseextractor.GermanStopList)

			chosenPhrases:= make(map[string]int)
			for _, phrase := range candidates {				
				_, ok :=mappedPhrases[phrase.Key]
				if ok {
					chosenPhrases[phrase.Key] = mappedPhrases[phrase.Key]
				}
			}	

			for _, phrase := range phraseextractor.RankByWordCount(chosenPhrases) {	
				fmt.Println(phrase.Key,phrase.Value)
			}
		}
				
		// if i > 0 {
		// 	endOfFile = true
		// }
		i = i + 1
	}
}





