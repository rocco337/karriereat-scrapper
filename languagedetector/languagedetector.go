package languagedetector

import "github.com/abadojack/whatlanggo"
import "karriereat-scrapper/filestorage"
import "karriereat-scrapper/dataaccess"
import  "encoding/json"
import  "fmt"

type LanguageDetector struct {
	options whatlanggo.Options
}

func (languageDetector *LanguageDetector) Init(){
	languageDetector.options= whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Eng: true,
			whatlanggo.Deu: true,
		},
	}
}

func (languageDetector *LanguageDetector) DetectAndSave(){
	fileReader :=new(filestorage.FileReader)
	fileReader.Init("result.json")
	
	line, endOfFile:=fileReader.ReadLine()
	for endOfFile==false{		
		jobsDetails :=new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		info := whatlanggo.DetectWithOptions(jobsDetails.Content, languageDetector.options)
		fmt.Println("Url: ", jobsDetails.Url, "Language:", info.Lang.Iso6391())
		line, endOfFile=fileReader.ReadLine()
	}
}