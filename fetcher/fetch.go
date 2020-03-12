package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rssFeed/db"
	"rssFeed/seed"
	"rssFeed/slack"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func getLastestStored(s seed.Seeder) string {
	query := fmt.Sprintf("select title from articles where feed_id='%s' and platform='%s' order by created_at desc limit 1", s.Identifier(), s.Platform())
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Printf("query articles with feedID %s error %+v", s.Identifier(), err)
		return ""
	}
	defer rows.Close()
	var a string
	for rows.Next() {
		if err := rows.Scan(&a); err != nil {
			log.Fatal(err)
		}
	}
	return a
}

func fetchArticles(s seed.Seeder, offset, limit int) ([]*db.Article, bool) {
	log.Printf("[%s] fetch %s", s.Name(), s.Link(offset, limit))
	resp, err := http.Get(s.Link(offset, limit))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonContent interface{}
	if err := json.Unmarshal(content, &jsonContent); err != nil {
		log.Fatal(err)
	}
	jc := jsonContent.(map[string]interface{})
	data := jc["data"].([]interface{})

	var articles []*db.Article
	for _, item := range data {
		v := item.(map[string]interface{})
		link := v["url"].(string)
		created := int64(v["created"].(float64))
		articles = append(articles, &db.Article{
			FeedName:    s.Name(),
			FeedID:      s.Identifier(),
			Platform:    s.Platform(),
			Title:       v["title"].(string),
			Link:        link,
			Description: v["excerpt"].(string),
			CreatedAt:   created,
		})
	}

	hasMore := false
	paging, ok := jc["paging"].(map[string]interface{})
	if ok {
		isEnd, ok := paging["is_end"].(bool)
		hasMore = ok && !isEnd
	}
	return articles, hasMore
}

func run(s seed.Seeder) {
	lastest := getLastestStored(s)
	log.Printf("get latest stored %s\n", lastest)

	offset := 0
	hasMore := true
	for hasMore {
		items, next := fetchArticles(s, offset, 10)
		if items != nil {
			offset += len(items)
		}
		hasMore = next

		index := len(items)
		for i, item := range items {
			if item.Title == lastest {
				index = i
				break
			}
		}
		if index > 0 {
			if err := db.SaveArticles(items[:index]); err != nil {
				log.Fatal(err)
			}
			// 可以利用消息系统来解耦
			for _, a := range items[:index] {
				// 发送slack消息
				slackMsg := fmt.Sprintf(`
					{
						"type": "mrkdwn",
						"text": "<%s|%s>"
					}
				`, a.Link, a.Title)
				slack.Notify(slackMsg)

				// 发送到pulisher服务
			}
		}
		if index != len(items) {
			return
		}
	}
}

func Sync() {
	run(seed.NewSeeder("zhihu", "prattle", "迷思"))
	run(seed.NewSeeder("zhihu", "milocode", "Milo的编程"))
}
