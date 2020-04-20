package httpServer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"rssFeed/db"
	"rssFeed/seed"
	"time"

	"github.com/gorilla/feeds"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slack-go/slack"
)

var signingSecret string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	signingSecret = os.Getenv("slackSigningSecret")
}

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

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Printf("slack.NewSecretsVerifier %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("slack.SlashCommandParse %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Printf("verifier.Ensure %+v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/keyword":
		params := &slack.Msg{Text: s.Text}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(s.Command))
		return
	}
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
