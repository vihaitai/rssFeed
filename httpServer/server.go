package httpServer

import (
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

//Start 开启服务
func Start(endpoint string) {
	log.Printf("[DEBUG]start http server on %s\n", endpoint)
	http.HandleFunc("/rss", rssListHandler)
	http.HandleFunc("/slash", slashCommandHandler)

	if err := http.ListenAndServe(endpoint, nil); err != nil {
		log.Fatal(err)
	}
}
