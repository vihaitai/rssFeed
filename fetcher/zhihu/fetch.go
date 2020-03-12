package zhihu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rssFeed/db"
	"rssFeed/slack"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const platform = "zhihu"

func getLastestStored(feedID string) string {
	query := fmt.Sprintf("select title from articles where feed_id='%s' and platform='%s' order by created_at desc limit 1", feedID, platform)
	rows, err := db.Conn.Query(query)
	if err != nil {
		log.Printf("query articles with feedID %s error %+v", feedID, err)
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

func fetchArticles(feedName, feedID, url string) ([]*db.Article, bool) {
	log.Printf("[%s] fetch %s", feedName, url)
	resp, err := http.Get(url)
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
			FeedName:    feedName,
			FeedID:      feedID,
			Platform:    platform,
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

func run(feedID, feedName string) {
	lastest := getLastestStored(feedID)
	log.Printf("get latest stored %s\n", lastest)

	offset := 0
	hasMore := true
	for hasMore {
		zhihuAPI := fmt.Sprintf("https://zhuanlan.zhihu.com/api/columns/%s/articles", feedID)
		articleURL := fmt.Sprintf("%s?offset=%d&limit=10", zhihuAPI, offset)
		items, next := fetchArticles(feedName, feedID, articleURL)
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
	run("prattle", "迷思")
	run("milocode", "Milo的编程")
}
