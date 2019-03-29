package main

import (
	"encoding/json"
	"fmt"
	"karriereat-scrapper/dataaccess"
	"karriereat-scrapper/filestorage"
	"karriereat-scrapper/pagescrappers"
	//"regexp"
	//"sort"
	"math"
	"strings"
	"time"
//"sync"
	//"github.com/bbalet/stopwords"
	//"github.com/dchest/stemmer/german"	
	"github.com/afjoseph/RAKE.Go"
	//"gopkg.in/jdkato/prose.v2"
	//"github.com/arpitgogia/rake"
	"github.com/goglue/tfidf"
	
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

	i:=0
	ds:=new(DocumentStorage)
	w := tfidf.New(ds)
	
	var r map[string]float64
	fullContent:=""
	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		//cleanContent := stopwords.CleanString(jobsDetails.Content, "de", true)	
		//cleanContent = strings.ToLower(cleanContent)
		// doc, err := prose.NewDocument(cleanContent)
		// if err != nil {
		// 	panic(err)
		// }

		// for _, tok := range doc.Tokens() {
		// 	//fmt.Println(tok.Text)	
		// 	fullContent = fullContent + " " + tok.Text		
		// }
		
		ds.documents = append(ds.documents, tfidf.StripSpacesLoop(jobsDetails.Content))
				
		//reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

		// cleanContent:=strings.ToLower(jobsDetails.Content)
		// cleanContent = reg.ReplaceAllString(cleanContent, " ")
		

		// words := strings.Fields(cleanContent)
		// sort.Sort(byLength(words))

		// stemmed :=""
		// for _, word := range words {
		// 	stemmed = stemmed + " " + german.Stemmer.Stem(word)
		// 	//fmt.Println(german.Stemmer.Stem(word))
		// }
		//fmt.Println(stemmed)
		//fullContent = fullContent + " " + stemmed

		fullContent= fullContent + " " + jobsDetails.Content		

		// }

		line, endOfFile = fileReader.ReadLine()
		endOfFile = true		
		if i >500{
			endOfFile = true		
		}
		i=i+1
	}

	// titles:=make([]string,0)
	// for _, document := range ds.documents{
	// 	dc := tfidf.StripSpacesLoop(document)	
	// 	r = w.Score(dc)

	// 	var keys []float64
	// 	inverted := make(map[float64]string)		
	// 	for key, value := range r {			
	// 		inverted[value]=key
	// 		keys = append(keys, value)			
	// 	}

	// 	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	// 	//fmt.Println("------------ Doc start ------------")
	// 	for _, k := range keys {
	// 		//fmt.Println(k, inverted[k])	
	// 		titles = append(titles,  inverted[k])
	// 		break	
	// 	}
	// 	//fmt.Println("------------ Doc end ------------")		
	// }
	// sort.Strings(titles)
	// for _, title :=range titles{
	// 	fmt.Println(title)
	// }
	
	candidates := rake.RunRakeI18N(fullContent, german)

	for _, candidate := range candidates {
		words :=strings.Split(candidate.Key, " ")
		r = ScoreTerms(w, ds, words)
			fmt.Println("------------ Term start ------------")
			for key, value := range r {
				fmt.Println(key, value)					
			}
		    fmt.Println("------------ Term end ------------")		
			//fmt.Printf("%s --> %f\n", candidate.Key, candidate.Value)
	}

	
	// rakeResult :=rake.WithText(fullContent)
	// for key, value := range rakeResult {
	// 	fmt.Println(key, value)
	// }
	
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


type DocumentStorage struct {
	documents []string
}

	// DocumentsWith receives a t term parameter and returns an unsigned integer of the documents count containing t
	func (s *DocumentStorage) DocumentsWith(t string) uint{
		count:=0

		for _, document := range s.documents{
			count =count + strings.Count(document, t)
			// if strings.Contains(document, t){
			// 	count = count + 1
			// }
		}
		return uint(count)
	}
	// Documents returns the total amount of documents within the storage
	func (s *DocumentStorage) Documents() uint{
		return uint(len(s.documents))
	}

	var german = []string{"ab", "aber", "ach", "acht", "achte", "achten", "achter", "achtes", "ag", "alle", "allein", "allem", "allen", "aller", "allerdings", "alles", "allgemeinen", "als", "also", "am", "an", "andere", "anderen", "andern", "anders", "au", "auch", "auf", "aus", "ausser", "außer", "ausserdem", "außerdem", "bald", "bei", "beide", "beiden", "beim", "beispiel", "bekannt", "bereits", "besonders", "besser", "besten", "bin", "bis", "bisher", "bist", "da", "dabei", "dadurch", "dafür", "dagegen", "daher", "dahin", "dahinter", "damals", "damit", "danach", "daneben", "dank", "dann", "daran", "darauf", "daraus", "darf", "darfst", "darin", "darüber", "darum", "darunter", "das", "dasein", "daselbst", "dass", "daß", "dasselbe", "davon", "davor", "dazu", "dazwischen", "dein", "deine", "deinem", "deiner", "dem", "dementsprechend", "demgegenüber", "demgemäss", "demgemäß", "demselben", "demzufolge", "den", "denen", "denn", "denselben", "der", "deren", "derjenige", "derjenigen", "dermassen", "dermaßen", "derselbe", "derselben", "des", "deshalb", "desselben", "dessen", "deswegen", "d.", "dich", "die", "diejenige", "diejenigen", "dies", "diese", "dieselbe", "dieselben", "diesem", "diesen", "dieser", "dieses", "dir", "doch", "dort", "drei", "drin", "dritte", "dritten", "dritter", "drittes", "du", "durch", "durchaus", "dürfen", "dürft", "durfte", "durften", "eben", "ebenso", "ehrlich", "ei", "eigen", "eigene", "eigenen", "eigener", "eigenes", "ein", "einander", "eine", "einem", "einen", "einer", "eines", "einige", "einigen", "einiger", "einiges", "einmal", "eins", "elf", "en", "ende", "endlich", "entweder", "er", "ernst", "erst", "erste", "ersten", "erster", "erstes", "es", "etwa", "etwas", "euch", "früher", "fünf", "fünfte", "fünften", "fünfter", "fünftes", "für", "gab", "ganz", "ganze", "ganzen", "ganzer", "ganzes", "gar", "gedurft", "gegen", "gegenüber", "gehabt", "gehen", "geht", "gekannt", "gekonnt", "gemacht", "gemocht", "gemusst", "genug", "gerade", "gern", "gesagt", "geschweige", "gewesen", "gewollt", "geworden", "gibt", "ging", "gleich", "gott", "gross", "groß", "grosse", "große", "grossen", "großen", "grosser", "großer", "grosses", "großes", "gut", "gute", "guter", "gutes", "habe", "haben", "habt", "hast", "hat", "hatte", "hätte", "hatten", "hätten", "heisst", "her", "heute", "hier", "hin", "hinter", "hoch", "ich", "ihm", "ihn", "ihnen", "ihr", "ihre", "ihrem", "ihren", "ihrer", "ihres", "im", "immer", "in", "indem", "infolgedessen", "ins", "irgend", "ist", "ja", "jahr", "jahre", "jahren", "je", "jede", "jedem", "jeden", "jeder", "jedermann", "jedermanns", "jedoch", "jemand", "jemandem", "jemanden", "jene", "jenem", "jenen", "jener", "jenes", "jetzt", "kam", "kann", "kannst", "kaum", "kein", "keine", "keinem", "keinen", "keiner", "kleine", "kleinen", "kleiner", "kleines", "kommen", "kommt", "können", "könnt", "konnte", "könnte", "konnten", "kurz", "lang", "lange", "leicht", "leide", "lieber", "los", "machen", "macht", "machte", "mag", "magst", "mahn", "man", "manche", "manchem", "manchen", "mancher", "manches", "mann", "mehr", "mein", "meine", "meinem", "meinen", "meiner", "meines", "mensch", "menschen", "mich", "mir", "mit", "mittel", "mochte", "möchte", "mochten", "mögen", "möglich", "mögt", "morgen", "muss", "muß", "müssen", "musst", "müsst", "musste", "mussten", "na", "nach", "nachdem", "nahm", "natürlich", "neben", "nein", "neue", "neuen", "neun", "neunte", "neunten", "neunter", "neuntes", "nicht", "nichts", "nie", "niemand", "niemandem", "niemanden", "noch", "nun", "nur", "ob", "oben", "oder", "offen", "oft", "ohne", "ordnung", "recht", "rechte", "rechten", "rechter", "rechtes", "richtig", "rund", "sa", "sache", "sagt", "sagte", "sah", "satt", "schlecht", "schluss", "schon", "sechs", "sechste", "sechsten", "sechster", "sechstes", "sehr", "sei", "seid", "seien", "sein", "seine", "seinem", "seinen", "seiner", "seines", "seit", "seitdem", "selbst", "sich", "sie", "sieben", "siebente", "siebenten", "siebenter", "siebentes", "sind", "so", "solang", "solche", "solchem", "solchen", "solcher", "solches", "soll", "sollen", "sollte", "sollten", "sondern", "sonst", "sowie", "später", "statt", "tag", "tage", "tagen", "tat", "teil", "tel", "tritt", "trotzdem", "tun", "über", "überhaupt", "übrigens", "uhr", "um", "und", "und?", "uns", "unser", "unsere", "unserer", "unter", "vergangenen", "viel", "viele", "vielem", "vielen", "vielleicht", "vier", "vierte", "vierten", "vierter", "viertes", "vom", "von", "vor", "wahr?", "während", "währenddem", "währenddessen", "wann", "war", "wäre", "waren", "wart", "warum", "was", "wegen", "weil", "weit", "weiter", "weitere", "weiteren", "weiteres", "welche", "welchem", "welchen", "welcher", "welches", "wem", "wen", "wenig", "wenige", "weniger", "weniges", "wenigstens", "wenn", "wer", "werde", "werden", "werdet", "wessen", "wie", "wieder", "will", "willst", "wir", "wird", "wirklich", "wirst", "wo", "wohl", "wollen", "wollt", "wollte", "wollten", "worden", "wurde", "würde", "wurden", "würden", "z.", "zehn", "zehnte", "zehnten", "zehnter", "zehntes", "zeit", "zu", "zuerst", "zugleich", "zum", "zunächst", "zur", "zurück", "zusammen", "zwanzig", "zwar", "zwei", "zweite", "zweiten", "zweiter", "zweites", "zwischen", "zwölf"}
	

	func ScoreTerms(w *tfidf.Weigher, ds *DocumentStorage, terms []string) map[string]float64 {	
		// tf terms frequencies within the given document
		tf := make(map[string]int)
		// tt total terms count within a given document
		tt := len(terms)
	
		for i := 0; i < tt; i++ {
			tf[terms[i]]++
		}
	
		// tft tf(t) term frequency of t
		tfidf := make(map[string]float64, len(tf))
	
		for term, freq := range tf {
			tft := float64(freq) / float64(tt)
			dwt := float64(ds.DocumentsWith(term))
	
			var idf float64
			if 0 == dwt {
				idf = 0
			} else {
				idf = math.Log10(
					float64(ds.Documents()) / dwt,
				)
			}
			tfidf[term] = tft * idf
		}
	
		return tfidf
	}