package httpServer

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"rssFeed/db"
	"rssFeed/seed"
	"time"

	"github.com/gorilla/feeds"
	_ "github.com/mattn/go-sqlite3"
)

func selectArticles(conn *sql.DB, s seed.Seeder) []db.Article {
	query := fmt.Sprintf("select id, title, link, description, created_at, feed_name from articles where platform=? and feed_id=? order by created_at desc")
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatalf("prepare %s error %+v", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(s.Platform(), s.Identifier())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var articles []db.Article
	for rows.Next() {
		var a = db.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Link, &a.Description, &a.CreatedAt, &a.FeedName)
		if err != nil {
			log.Printf("article scan error %+v", err)
			continue
		}
		articles = append(articles, a)
	}
	return articles
}

func rssListHandler(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	platform := query["platform"][0]
	format := query["format"][0]
	feedID := query["feed"][0]
	s := seed.NewSeeder(platform, feedID, "")
	feed := &feeds.Feed{
		Link: &feeds.Link{Href: s.Home()},
	}
	var feedItems []*feeds.Item
	articles := selectArticles(db.Conn, s)
	for _, a := range articles {
		if feed.Title == "" {
			feed.Title = a.FeedName
		}
		feedItems = append(feedItems, &feeds.Item{
			Id:          a.ID,
			Title:       a.Title,
			Link:        &feeds.Link{Href: a.Link},
			Description: a.Description,
			Created:     time.Unix(a.CreatedAt, 0),
		})
	}
	feed.Items = feedItems

	var (
		payload string
		err     error
	)
	switch format {
	case "rss":
		payload, err = feed.ToRss()
	case "json":
		payload, err = feed.ToJSON()
	}

	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(writer, payload)
}
