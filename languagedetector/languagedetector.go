package languagedetector

import "github.com/abadojack/whatlanggo"
import "karriereat-scrapper/filestorage"
import "karriereat-scrapper/dataaccess"
import "encoding/json"

type LanguageDetector struct {
	options        whatlanggo.Options
	jobsDataAccess *dataaccess.JobsDataAccess
}

func (languageDetector *LanguageDetector) Init() {
	languageDetector.options = whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Eng: true,
			whatlanggo.Deu: true,
		},
	}

	languageDetector.jobsDataAccess = new(dataaccess.JobsDataAccess)
	languageDetector.jobsDataAccess.Init()
}

func (languageDetector *LanguageDetector) Detect(content string) string {
	info := whatlanggo.DetectWithOptions(content, languageDetector.options)
	return info.Lang.Iso6391()
}

func (languageDetector *LanguageDetector) DetectAndSave() {
	fileReader := new(filestorage.FileReader)
	fileReader.Init("result.json")

	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		info := whatlanggo.DetectWithOptions(jobsDetails.Content, languageDetector.options)
		languageDetector.jobsDataAccess.SetLanguage(jobsDetails.Url, info.Lang.Iso6391())
		line, endOfFile = fileReader.ReadLine()
	}
}
