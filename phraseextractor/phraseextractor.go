package phraseextractor

import "karriereat-scrapper/filestorage"
import (
	"encoding/json"
	"fmt"
	"karriereat-scrapper/dataaccess"
	"karriereat-scrapper/languagedetector"
    "github.com/dchest/stemmer/german"
	"sync"
	"sort"
	rake "github.com/afjoseph/RAKE.Go"
)

type Phraseextractor struct {

}

func (extractor *Phraseextractor) MapPhrases(content string, mappedPhrases map[string]int, mutex *sync.Mutex, wg *sync.WaitGroup, counter int) {
	defer wg.Done()
	defer fmt.Println(counter)

	ger := german.Stemmer			
	stemmed:= ger.Stem(content)

	candidates := rake.RunRakeI18N(stemmed, GermanStopList)
	mutex.Lock()
	for _, phrase := range candidates {	
		mappedPhrases[phrase.Key]++
	}	
	mutex.Unlock()			
}

func (extractor *Phraseextractor) ExtractPhrasesFromFile(fileReader *filestorage.FileReader,languagedetector *languagedetector.LanguageDetector) map[string]int{
	i := 0
	mappedPhrases:=make(map[string]int)
	var mutex = &sync.Mutex{}
	var myWaitGroup sync.WaitGroup

	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		jobsDetails := new(dataaccess.JobsDetails)
		if err := json.Unmarshal([]byte(line), &jobsDetails); err != nil {
			panic(err)
		}

		lang := languagedetector.Detect(jobsDetails.Content)
		if lang == "de" {
			myWaitGroup.Add(2) 
			go extractor.MapPhrases(jobsDetails.Content, mappedPhrases, mutex, &myWaitGroup, i)		
			go extractor.MapPhrases(jobsDetails.Title, mappedPhrases, mutex, &myWaitGroup, i)			
		}

		line, endOfFile = fileReader.ReadLine()
		// if i > 100 {
		// 	endOfFile = true
		// }
		i = i + 1
	}
	myWaitGroup.Wait()

	for key, _:= range mappedPhrases{
		if len(key)<=1{
			delete(mappedPhrases, key);
		}else{
			for _, stopWord := range GermanStopList {
				if stopWord == key {
					delete(mappedPhrases, key);
					break;
				}
			}
		}
	}
	return mappedPhrases
}

func (extractor *Phraseextractor) WritePhrasesToFile(phrases map[string]int,fileStorage *filestorage.FileStorage){
	for _, pair := range RankByWordCount(phrases) {		
		json, err := json.Marshal(pair)
			if err != nil {
				panic(err)
			}

		fileStorage.AppendLine(string(json))
	}
}

func (extractor *Phraseextractor) ReadPhrasesFromFile(fileReader *filestorage.FileReader) map[string]int{
	mappedPhrases:=make(map[string]int)
	line, endOfFile := fileReader.ReadLine()
	for endOfFile == false {
		phrase := new(Pair)
		if err := json.Unmarshal([]byte(line), &phrase); err != nil {
			panic(err)
		}

		mappedPhrases[phrase.Key] = int(phrase.Value)
		line, endOfFile = fileReader.ReadLine()
	}

	return mappedPhrases
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

func RankByWordCount(wordFrequencies map[string]int) PairList{
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
	  pl[i] = Pair{k, float64(v)}
	  i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
  }

  var GermanStopList = []string{"ab", "aber", "ach", "acht", "achte", "achten", "achter", "achtes", "ag", "alle", "allein", "allem", "allen", "aller", "allerdings", "alles", "allgemeinen", "als", "also", "am", "an", "andere", "anderen", "andern", "anders", "au", "auch", "auf", "aus", "ausser", "außer", "ausserdem", "außerdem", "bald", "bei", "beide", "beiden", "beim", "beispiel", "bekannt", "bereits", "besonders", "besser", "besten", "bin", "bis", "bisher", "bist", "da", "dabei", "dadurch", "dafür", "dagegen", "daher", "dahin", "dahinter", "damals", "damit", "danach", "daneben", "dank", "dann", "daran", "darauf", "daraus", "darf", "darfst", "darin", "darüber", "darum", "darunter", "das", "dasein", "daselbst", "dass", "daß", "dasselbe", "davon", "davor", "dazu", "dazwischen", "dein", "deine", "deinem", "deiner", "dem", "dementsprechend", "demgegenüber", "demgemäss", "demgemäß", "demselben", "demzufolge", "den", "denen", "denn", "denselben", "der", "deren", "derjenige", "derjenigen", "dermassen", "dermaßen", "derselbe", "derselben", "des", "deshalb", "desselben", "dessen", "deswegen", "d.", "dich", "die", "diejenige", "diejenigen", "dies", "diese", "dieselbe", "dieselben", "diesem", "diesen", "dieser", "dieses", "dir", "doch", "dort", "drei", "drin", "dritte", "dritten", "dritter", "drittes", "du", "durch", "durchaus", "dürfen", "dürft", "durfte", "durften", "eben", "ebenso", "ehrlich", "ei", "eigen", "eigene", "eigenen", "eigener", "eigenes", "ein", "einander", "eine", "einem", "einen", "einer", "eines", "einige", "einigen", "einiger", "einiges", "einmal", "eins", "elf", "en", "ende", "endlich", "entweder", "er", "ernst", "erst", "erste", "ersten", "erster", "erstes", "es", "etwa", "etwas", "euch", "früher", "fünf", "fünfte", "fünften", "fünfter", "fünftes", "für", "gab", "ganz", "ganze", "ganzen", "ganzer", "ganzes", "gar", "gedurft", "gegen", "gegenüber", "gehabt", "gehen", "geht", "gekannt", "gekonnt", "gemacht", "gemocht", "gemusst", "genug", "gerade", "gern", "gesagt", "geschweige", "gewesen", "gewollt", "geworden", "gibt", "ging", "gleich", "gott", "gross", "groß", "grosse", "große", "grossen", "großen", "grosser", "großer", "grosses", "großes", "gut", "gute", "guter", "gutes", "habe", "haben", "habt", "hast", "hat", "hatte", "hätte", "hatten", "hätten", "heisst", "her", "heute", "hier", "hin", "hinter", "hoch", "ich", "ihm", "ihn", "ihnen", "ihr", "ihre", "ihrem", "ihren", "ihrer", "ihres", "im", "immer", "in", "indem", "infolgedessen", "ins", "irgend", "ist", "ja", "jahr", "jahre", "jahren", "je", "jede", "jedem", "jeden", "jeder", "jedermann", "jedermanns", "jedoch", "jemand", "jemandem", "jemanden", "jene", "jenem", "jenen", "jener", "jenes", "jetzt", "kam", "kann", "kannst", "kaum", "kein", "keine", "keinem", "keinen", "keiner", "kleine", "kleinen", "kleiner", "kleines", "kommen", "kommt", "können", "könnt", "konnte", "könnte", "konnten", "kurz", "lang", "lange", "leicht", "leide", "lieber", "los", "machen", "macht", "machte", "mag", "magst", "mahn", "man", "manche", "manchem", "manchen", "mancher", "manches", "mann", "mehr", "mein", "meine", "meinem", "meinen", "meiner", "meines", "mensch", "menschen", "mich", "mir", "mit", "mittel", "mochte", "möchte", "mochten", "mögen", "möglich", "mögt", "morgen", "muss", "muß", "müssen", "musst", "müsst", "musste", "mussten", "na", "nach", "nachdem", "nahm", "natürlich", "neben", "nein", "neue", "neuen", "neun", "neunte", "neunten", "neunter", "neuntes", "nicht", "nichts", "nie", "niemand", "niemandem", "niemanden", "noch", "nun", "nur", "ob", "oben", "oder", "offen", "oft", "ohne", "ordnung", "recht", "rechte", "rechten", "rechter", "rechtes", "richtig", "rund", "sa", "sache", "sagt", "sagte", "sah", "satt", "schlecht", "schluss", "schon", "sechs", "sechste", "sechsten", "sechster", "sechstes", "sehr", "sei", "seid", "seien", "sein", "seine", "seinem", "seinen", "seiner", "seines", "seit", "seitdem", "selbst", "sich", "sie", "sieben", "siebente", "siebenten", "siebenter", "siebentes", "sind", "so", "solang", "solche", "solchem", "solchen", "solcher", "solches", "soll", "sollen", "sollte", "sollten", "sondern", "sonst", "sowie", "später", "statt", "tag", "tage", "tagen", "tat", "teil", "tel", "tritt", "trotzdem", "tun", "über", "überhaupt", "übrigens", "uhr", "um", "und", "und?", "uns", "unser", "unsere", "unserer", "unter", "vergangenen", "viel", "viele", "vielem", "vielen", "vielleicht", "vier", "vierte", "vierten", "vierter", "viertes", "vom", "von", "vor", "wahr?", "während", "währenddem", "währenddessen", "wann", "war", "wäre", "waren", "wart", "warum", "was", "wegen", "weil", "weit", "weiter", "weitere", "weiteren", "weiteres", "welche", "welchem", "welchen", "welcher", "welches", "wem", "wen", "wenig", "wenige", "weniger", "weniges", "wenigstens", "wenn", "wer", "werde", "werden", "werdet", "wessen", "wie", "wieder", "will", "willst", "wir", "wird", "wirklich", "wirst", "wo", "wohl", "wollen", "wollt", "wollte", "wollten", "worden", "wurde", "würde", "wurden", "würden", "z.", "zehn", "zehnte", "zehnten", "zehnter", "zehntes", "zeit", "zu", "zuerst", "zugleich", "zum", "zunächst", "zur", "zurück", "zusammen", "zwanzig", "zwar", "zwei", "zweite", "zweiten", "zweiter", "zweites", "zwischen", "zwölf"}