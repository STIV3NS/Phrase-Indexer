package main

import (
    "fmt"
    "net/http"
    "regexp"
    "strings"
    "github.com/PuerkitoBio/goquery"
)

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

func normalize(text *string) {
    nonWord, _ := regexp.Compile("[0-9`~!@#$%^?&*()_+-=\\[\\]{}|'\";:/.,><]")
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

func filter(src *[]phraseCnt, predicate func(phraseCnt) bool) *[]phraseCnt {
    filtered := make([]phraseCnt, 0)
    
    for _, elem := range *src {
        if predicate(elem) {
            filtered = append(filtered, elem)
        }
    }
    
    return &filtered
}
