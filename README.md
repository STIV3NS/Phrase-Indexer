# Phrase-Indexer
Tool for creating a ranking of phrases (under a given selector) in a forum thread

#### Build

```bash
$ make
```

#### Usage
```bash
$ ./phrase_indexer --help
Usage of ./phrase_indexer:
  -threadURL string
      [REQUIRED] URL to threadURL that is meant to be indexed
  -selector string
    	[REQUIRED] Selector for searching for interesting parts of the document
  -start uint
    	[OPTIONAL] Page number on which to start indexing (default 1)
  -end uint
    	[REQUIRED] Page number on which to end indexing
  -exclude string
    	[OPTIONAL] Path to file that contains phrases to exclude from output
                 [text file, whitespace separated]
  -workers uint
    	[OPTIONAL] Number of workers involved to parsing thread sites (default 100)
  -limit int
    	[OPTIONAL] Limit output to top #{value} entries (default 2147483647)

```

#### Example
```bash
$ ./phrase_indexer \
-threadURL="https://4programmers.net/Forum/Off-Topic/141606-programistyczne_wtf_jakie_was_spotkaly?page=" \
-selector=".online, .offline" -end=100 -limit=10

218 			rnd
128 			somekind
100 			azarien
94 			koziolek
66 			marekr
54 			marooned
42 			krolik
42 			monk
42 			demonical
42 			wibowit


```

#### Dependencies

[goquery](https://github.com/PuerkitoBio/goquery)
