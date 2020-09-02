package main

import (
    "io/ioutil"
    "strings"
)

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