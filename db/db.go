package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

type Article struct {
	ID          string
	Title       string
	Link        string
	Description string
	Platform    string
	FeedName    string
	FeedID      string
	CreatedAt   int64
}

func init() {
	conn, err := sql.Open("sqlite3", "/home/ubuntu/rssFeed/rss.db")
	if err != nil {
		log.Fatal("can not open rss.db")
	}
	Conn = conn
	_, err = Conn.Exec(`
		CREATE TABLE IF NOT EXISTS articles(
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			link TEXT NOT NULL UNIQUE,
			feed_name TEXT NOT NULL,
			feed_id TEXT NOT NULL,
			platform TEXT NOT NULL,
			description TEXT NOT NULL,
			created_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		log.Printf("create articles table error %+v", err)
	}
}

func SaveArticles(items []*Article) error {
	query_tpl := "insert into articles(`feed_name`, `feed_id`, `platform`,  `title`, `link`, `description`, `created_at`) values %s"
	values := make([]string, len(items))
	for i, item := range items {
		values[i] = fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%d')", item.FeedName, item.FeedID, item.Platform, item.Title, item.Link, item.Description, item.CreatedAt)
	}
	query := fmt.Sprintf(query_tpl, strings.Join(values, ","))
	result, err := Conn.Exec(query)
	if err != nil {
		log.Printf("exec %s return %+v\n", query, result)
	}
	return err
}
