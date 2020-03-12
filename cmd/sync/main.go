package main

import (
	"rssFeed/fetcher"
)

func main() {
	//var source string
	//flag.StringVar(&source, "s", "", "")
	//flag.Parse()
	fetcher.Sync()
}
