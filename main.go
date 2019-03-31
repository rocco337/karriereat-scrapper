package main

import (
	"encoding/json"
	"fmt"
	"karriereat-scrapper/dataaccess"
	"karriereat-scrapper/filestorage"
	"karriereat-scrapper/languagedetector"
	"karriereat-scrapper/pagescrappers"
	"sort"
	"strings"
	"time"

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

	i := 0
	fullContent := ""
	files := make([]string, 0)
	titles := make([]string, 0)

	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		lang := languagedetector.Detect(jobsDetails.Content)
		if lang == "de" {
			fullContent += " " + jobsDetails.Content
			files = append(files, strings.ToLower(jobsDetails.Content))
			titles = append(titles, strings.ToLower(jobsDetails.Title))
		}

		line, endOfFile = fileReader.ReadLine()
		// if i > 100 {
		// 	endOfFile = true
		// }
		i = i + 1
	}

	phrases := make(map[string]float64, 0)
	for _, file := range files {
		candidates := rake.RunRakeI18N(file, german)

		for _, candidate := range candidates {
			score := float64(candidate.Value) / 100
			count, _ := AppearsInHowManyFiles(titles, candidate.Key)
			score += float64(count)

			count, _ = AppearsInHowManyFiles(files, candidate.Key)
			score += float64(count) / float64(len(files))

			phrases[candidate.Key] = score
		}

	}

	SortAndPrintResults(phrases)

	elapsed := time.Since(start)
	fmt.Println("Elapsed", elapsed)
}

func SortAndPrintResults(phrases map[string]float64) {
	p := make(PairList, len(phrases))

	i := 0
	for k, v := range phrases {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(p)

	fmt.Printf("Post-sorted: ")
	for _, k := range p {
		fmt.Println(k.Value, k.Key)
	}
	fmt.Println("")
}

func AppearsInHowManyFiles(files []string, phrase string) (int, int) {
	phrase = strings.ToLower(phrase)
	filesContaingPhrase := 0
	totalNumberOfOccurences := 0
	for _, file := range files {
		if strings.Contains(file, phrase) {
			filesContaingPhrase++
		}

		totalNumberOfOccurences += strings.Count(file, phrase)
	}

	return filesContaingPhrase, totalNumberOfOccurences
}

type Pair struct {
	Key   string
	Value float64
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

var german = []string{"ab", "aber", "ach", "acht", "achte", "achten", "achter", "achtes", "ag", "alle", "allein", "allem", "allen", "aller", "allerdings", "alles", "allgemeinen", "als", "also", "am", "an", "andere", "anderen", "andern", "anders", "au", "auch", "auf", "aus", "ausser", "außer", "ausserdem", "außerdem", "bald", "bei", "beide", "beiden", "beim", "beispiel", "bekannt", "bereits", "besonders", "besser", "besten", "bin", "bis", "bisher", "bist", "da", "dabei", "dadurch", "dafür", "dagegen", "daher", "dahin", "dahinter", "damals", "damit", "danach", "daneben", "dank", "dann", "daran", "darauf", "daraus", "darf", "darfst", "darin", "darüber", "darum", "darunter", "das", "dasein", "daselbst", "dass", "daß", "dasselbe", "davon", "davor", "dazu", "dazwischen", "dein", "deine", "deinem", "deiner", "dem", "dementsprechend", "demgegenüber", "demgemäss", "demgemäß", "demselben", "demzufolge", "den", "denen", "denn", "denselben", "der", "deren", "derjenige", "derjenigen", "dermassen", "dermaßen", "derselbe", "derselben", "des", "deshalb", "desselben", "dessen", "deswegen", "d.", "dich", "die", "diejenige", "diejenigen", "dies", "diese", "dieselbe", "dieselben", "diesem", "diesen", "dieser", "dieses", "dir", "doch", "dort", "drei", "drin", "dritte", "dritten", "dritter", "drittes", "du", "durch", "durchaus", "dürfen", "dürft", "durfte", "durften", "eben", "ebenso", "ehrlich", "ei", "eigen", "eigene", "eigenen", "eigener", "eigenes", "ein", "einander", "eine", "einem", "einen", "einer", "eines", "einige", "einigen", "einiger", "einiges", "einmal", "eins", "elf", "en", "ende", "endlich", "entweder", "er", "ernst", "erst", "erste", "ersten", "erster", "erstes", "es", "etwa", "etwas", "euch", "früher", "fünf", "fünfte", "fünften", "fünfter", "fünftes", "für", "gab", "ganz", "ganze", "ganzen", "ganzer", "ganzes", "gar", "gedurft", "gegen", "gegenüber", "gehabt", "gehen", "geht", "gekannt", "gekonnt", "gemacht", "gemocht", "gemusst", "genug", "gerade", "gern", "gesagt", "geschweige", "gewesen", "gewollt", "geworden", "gibt", "ging", "gleich", "gott", "gross", "groß", "grosse", "große", "grossen", "großen", "grosser", "großer", "grosses", "großes", "gut", "gute", "guter", "gutes", "habe", "haben", "habt", "hast", "hat", "hatte", "hätte", "hatten", "hätten", "heisst", "her", "heute", "hier", "hin", "hinter", "hoch", "ich", "ihm", "ihn", "ihnen", "ihr", "ihre", "ihrem", "ihren", "ihrer", "ihres", "im", "immer", "in", "indem", "infolgedessen", "ins", "irgend", "ist", "ja", "jahr", "jahre", "jahren", "je", "jede", "jedem", "jeden", "jeder", "jedermann", "jedermanns", "jedoch", "jemand", "jemandem", "jemanden", "jene", "jenem", "jenen", "jener", "jenes", "jetzt", "kam", "kann", "kannst", "kaum", "kein", "keine", "keinem", "keinen", "keiner", "kleine", "kleinen", "kleiner", "kleines", "kommen", "kommt", "können", "könnt", "konnte", "könnte", "konnten", "kurz", "lang", "lange", "leicht", "leide", "lieber", "los", "machen", "macht", "machte", "mag", "magst", "mahn", "man", "manche", "manchem", "manchen", "mancher", "manches", "mann", "mehr", "mein", "meine", "meinem", "meinen", "meiner", "meines", "mensch", "menschen", "mich", "mir", "mit", "mittel", "mochte", "möchte", "mochten", "mögen", "möglich", "mögt", "morgen", "muss", "muß", "müssen", "musst", "müsst", "musste", "mussten", "na", "nach", "nachdem", "nahm", "natürlich", "neben", "nein", "neue", "neuen", "neun", "neunte", "neunten", "neunter", "neuntes", "nicht", "nichts", "nie", "niemand", "niemandem", "niemanden", "noch", "nun", "nur", "ob", "oben", "oder", "offen", "oft", "ohne", "ordnung", "recht", "rechte", "rechten", "rechter", "rechtes", "richtig", "rund", "sa", "sache", "sagt", "sagte", "sah", "satt", "schlecht", "schluss", "schon", "sechs", "sechste", "sechsten", "sechster", "sechstes", "sehr", "sei", "seid", "seien", "sein", "seine", "seinem", "seinen", "seiner", "seines", "seit", "seitdem", "selbst", "sich", "sie", "sieben", "siebente", "siebenten", "siebenter", "siebentes", "sind", "so", "solang", "solche", "solchem", "solchen", "solcher", "solches", "soll", "sollen", "sollte", "sollten", "sondern", "sonst", "sowie", "später", "statt", "tag", "tage", "tagen", "tat", "teil", "tel", "tritt", "trotzdem", "tun", "über", "überhaupt", "übrigens", "uhr", "um", "und", "und?", "uns", "unser", "unsere", "unserer", "unter", "vergangenen", "viel", "viele", "vielem", "vielen", "vielleicht", "vier", "vierte", "vierten", "vierter", "viertes", "vom", "von", "vor", "wahr?", "während", "währenddem", "währenddessen", "wann", "war", "wäre", "waren", "wart", "warum", "was", "wegen", "weil", "weit", "weiter", "weitere", "weiteren", "weiteres", "welche", "welchem", "welchen", "welcher", "welches", "wem", "wen", "wenig", "wenige", "weniger", "weniges", "wenigstens", "wenn", "wer", "werde", "werden", "werdet", "wessen", "wie", "wieder", "will", "willst", "wir", "wird", "wirklich", "wirst", "wo", "wohl", "wollen", "wollt", "wollte", "wollten", "worden", "wurde", "würde", "wurden", "würden", "z.", "zehn", "zehnte", "zehnten", "zehnter", "zehntes", "zeit", "zu", "zuerst", "zugleich", "zum", "zunächst", "zur", "zurück", "zusammen", "zwanzig", "zwar", "zwei", "zweite", "zweiten", "zweiter", "zweites", "zwischen", "zwölf"}
