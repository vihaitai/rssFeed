package main

import (
	"flag"
	"rssFeed/fetcher/zhihu"
)

func main() {
	var source string
	flag.StringVar(&source, "s", "", "")
	flag.Parse()

	switch source {
	case "zhihu":
		zhihu.Sync()
	default:
		println("unknow souce " + source)
	}
}
