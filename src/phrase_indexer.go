package main

import (
    "flag"
    "fmt"
    "log"
    "math"
    "os"
    "sort"
    "strings"
    "sync"
    "github.com/PuerkitoBio/goquery"
)

type phraseCnt struct {
    phrase string
    count  uint32
}

type arguments struct {
    threadUrl string
    selector  string
    exclude   string
    startAt   uint
    endAt     uint
    limit     int
    nworkers  uint
}

func getArguments() arguments {
    const sREQUIRED = ""
    const iREQUIRED = 0
    const DEFAULT_WORKERS = 100
    
    var threadUrl, selector, exclude string
    var startAt, endAt, workers uint
    var limit int
    
    flag.StringVar(&threadUrl, "threadUrl", sREQUIRED,
                    "[REQUIRED] URL to thread that is meant to be indexed")
    flag.StringVar(&selector, "selector", sREQUIRED,
                    "[REQUIRED] Selector for searching for interesting parts of the document")
    flag.UintVar(&startAt, "startAt", 1,
                    "[OPTIONAL] Page number on which to start indexing")
    flag.UintVar(&endAt, "endAt", iREQUIRED,
                    "[REQUIRED] Page number on which to end indexing")

    flag.StringVar(&exclude, "exclude", "",
                    "[OPTIONAL] Path to file that contains phrases to exclude from output")
    flag.IntVar(&limit, "limit", math.MaxInt32,
                    "[OPTIONAL] Limit output to top #{value} entries")
    flag.UintVar(&workers, "workers", DEFAULT_WORKERS,
                    "[OPTIONAL] Number of workers involved to parsing thread sites")

    flag.Parse()
    
    if endAt < startAt || threadUrl == "" || selector == ""  {
        fmt.Fprintf(os.Stderr, "Missing arguments; --help for more information\n")
        os.Exit(1)
    }

    info, err2 := os.Stat(exclude)
    if strings.Compare(exclude, "") != 0 && (os.IsNotExist(err2) || info.IsDir()) {
        fmt.Fprintf(os.Stderr, "Path: %v does not exist or is a directory.\n", exclude)
        os.Exit(1)
    }

    jobSize := endAt - startAt + 1
    if workers > jobSize {
        workers = jobSize
    }

    return arguments{threadUrl, selector, exclude, startAt, endAt, limit, workers}
}

func main() {
    args := getArguments()
    
    collectorChan, result := spawnCollector(args.nworkers)
    jobs, wg := spawnWorkers(args.selector, args.nworkers, collectorChan)
    initJobs(jobs, args.threadUrl, args.startAt, args.endAt)
    
    wg.Wait()
    close(collectorChan)
    
    printOutRanking( <-result, args.limit, args.exclude )
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

func spawnCollector(bufferSize uint) (collectorChan chan *map[string]uint32, result chan *map[string]uint32) {
    collectorChan = make(chan *map[string]uint32, bufferSize)
    result = make(chan *map[string]uint32)
    
    go collector(collectorChan, result)
    
    return
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
        
        bufferSize := howMany
        jobs = make(chan string, bufferSize)
        
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
        jobs <- fmt.Sprintf("%v%v", threadURL, i)
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
            selection.Remove()
            rawText := selection.Text()
            normalize(&rawText)
            
            phrases := strings.Fields(rawText)
            
            for _, phrase := range phrases {
                phrasesCount[phrase]++
            }
        })
        
        err := resp.Body.Close()
        if err != nil {
            panic(err)
        }
        
        collector <- &phrasesCount
        log.Printf("Job done: %v\n", job)
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
