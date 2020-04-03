package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sort"
	"log"
	"io/ioutil"
)

type phraseCnt struct {
	phrase string
	count  uint32
}

func main() {
	threadURL, selector, exclude, start, end, nworkers, limit := getArguments()

	collectorChan := make(chan *map[string]uint32)
	result := make(chan *map[string]uint32)

	go collector(collectorChan, result)

	jobs, wg := spawnWorkers(selector, nworkers, collectorChan)
	initJobs(jobs, threadURL, start, end)



	wg.Wait()
	close(collectorChan)

	printOutRanking( <-result, limit, exclude )

}

func printOutRanking(phraseCounts *map[string]uint32, limit int, exclude string) {
	ranking := sortByPhraseCount(phraseCounts)

	if strings.Compare(exclude, "") != 0 {
		applyExclusions(&ranking, exclude)
	}

	for i, elem := range *ranking {
		if i >= limit {
			break
		}

		fmt.Printf("%v \t\t\t% v\n", elem.count, elem.phrase)
	}
}

func collector(input <-chan *map[string]uint32, result chan<- *map[string]uint32) {
	state := make(map[string]uint32)

	for iteration := range input {
		for phrase, count := range *iteration {
			state[phrase] += count
		}
	}

	result <- &state
}

func spawnWorkers(selector string, howMany uint, collector chan<- *map[string]uint32) (
	jobs chan string, wg *sync.WaitGroup) {
	jobs = make(chan string)

	var _wg sync.WaitGroup
	wg = &_wg

	var i uint
	for i = 0; i < howMany; i++ {
		wg.Add(1)
		go worker(selector, jobs, collector, wg)
	}

	return
}

func initJobs(jobs chan<- string, threadURL string, start uint, end uint) {
	for i := start; i <= end; i++ {
		jobs <- threadURL + fmt.Sprintf("%v", i)
	}
	close(jobs)
}

func worker(selector string, jobs <-chan string, collector chan<- *map[string]uint32, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		resp := getResponse(job)
		doc := *getHtml(resp)

		phrasesCount := make(map[string]uint32)

		doc.Find(selector).Each(func(i int, selection *goquery.Selection){

			rawText := selection.Text()
			normalize(&rawText)

			phrases := strings.Fields(rawText)

			for _, phrase := range phrases {
				phrasesCount[phrase]++
			}

		})

		resp.Body.Close()

		collector <- &phrasesCount
		log.Printf("Job done: %v\n", job)
	}
}

func getResponse(site string) *http.Response {
	resp, err := http.Get(site)

	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Response from %v has status code %v. Aborting.", site, resp.StatusCode))
	}

	return resp
}

func getHtml(resp *http.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	return doc
}

func getArguments() (threadURL, selector, exclude string, start, end, workers uint, limit int) {
	const sREQUIRED = ""
	const iREQUIRED = 0
	const DEFAULT_WORKERS = 100

	flag.StringVar(&threadURL, "threadURL", sREQUIRED,
		"[REQUIRED] URL to threadURL that is meant to be indexed")
	flag.StringVar(&selector, "selector", sREQUIRED,
		"[REQUIRED] Selector for searching for interesting parts of the document")
	flag.UintVar(&start, "start", 1,
		"[OPTIONAL] Page number on which to start indexing")
	flag.UintVar(&end, "end", iREQUIRED,
		"[REQUIRED] Page number on which to end indexing")

	flag.StringVar(&exclude, "exclude", "",
		"[OPTIONAL] Path to file that contains phrases to exclude from output")
	flag.IntVar(&limit, "limit", math.MaxInt32,
		"[OPTIONAL] Limit output to top #{value} entries")
	flag.UintVar(&workers, "workers", DEFAULT_WORKERS,
		"[OPTIONAL] Number of workers involved to parsing thread sites")

	flag.Parse()

	if end == 0 || threadURL == "" || selector == ""  {
		fmt.Fprintf(os.Stderr, "Missing arguments; --help for more information\n")
		os.Exit(1)
	}

	info, err2 := os.Stat(exclude)
	if strings.Compare(exclude, "") != 0 && (os.IsNotExist(err2) || info.IsDir()) {
		fmt.Fprintf(os.Stderr, "Path: %v does not exist or is a directory.\n", exclude)
		os.Exit(1)
	}

	jobSize := end - start + 1
	if workers > jobSize {
		workers = jobSize
	}

	return
}

func normalize(text *string) {
	nonWord, _ := regexp.Compile("[0-9`~!@#$%^&*()_+-=\\[\\]{}|'\";:/.,><]")
	*text = nonWord.ReplaceAllLiteralString(*text, "")
	*text = strings.ToLower(*text)
	replacePolishDiacritics(text)
}

func replacePolishDiacritics(text *string) {
	a, _ := regexp.Compile("ą")
	c, _ := regexp.Compile("ć")
	e, _ := regexp.Compile("ę")
	l, _ := regexp.Compile("ł")
	o, _ := regexp.Compile("ó")
	s, _ := regexp.Compile("ś")
	z, _ := regexp.Compile("[żź]")

	replacements := map[*regexp.Regexp]string{
		a : "a",
		c : "c",
		e : "e",
		l : "l",
		o : "o",
		s : "s",
		z : "z",
	}

	for regexptr, repl := range replacements {
		*text = regexptr.ReplaceAllLiteralString(*text, repl)
	}
}

func sortByPhraseCount(phraseCounts *map[string]uint32) *[]phraseCnt {
	ranking := make([]phraseCnt, len(*phraseCounts))

	var i uint32 = 0
	for phrase, cnt := range *phraseCounts {
		ranking[i].phrase = phrase
		ranking[i].count = cnt

		i++
	}

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].count > ranking[j].count
	})

	return &ranking
}

func filter(src *[]phraseCnt, predicate func(phraseCnt) bool) *[]phraseCnt {
	filtered := make([]phraseCnt, 0)

	for _, elem := range *src {
		if predicate(elem) {
			filtered = append(filtered, elem)
		}
	}

	return &filtered
}

func applyExclusions(ranking **[]phraseCnt, exclude string) {
	exclusions := getExclusions(exclude)

	*ranking = filter(*ranking, func(elem phraseCnt) bool {
		_, present := exclusions[elem.phrase]
		return !present
	})
}

func getExclusions(filePath string) (exclusions map[string]bool) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	exclusions = make(map[string]bool)

	str := string(data)
	normalize(&str)

	for _, str := range strings.Fields(str) {
		exclusions[str] = true
	}

	return
}
