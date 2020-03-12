package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"rssFeed/db"
	"rssFeed/pubServer"
	"time"
)

func selectArticles(conn *sql.DB, source string) []db.Article {
	query := fmt.Sprintf("select id, title, link, description, created_at, source from articles where source=? order by created_at")
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatalf("prepare %s error %+v", query, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(source)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var articles []db.Article
	for rows.Next() {
		var a = db.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Link, &a.Description, &a.CreatedAt, &a.Source)
		if err != nil {
			log.Printf("article scan error %+v", err)
			continue
		}
		articles = append(articles, a)
	}
	return articles
}

func main() {
	client, err := pubServer.NewREQClient("tcp://127.0.0.1:8002")
	if err != nil {
		panic(err)
	}
	articles := selectArticles(db.Conn, "迷思")
	total := 0
	for _, article := range articles {
		msg, err := json.Marshal(article)
		if err != nil {
			panic(err)
		}
		client.SendREQMsg(msg)
		time.Sleep(500 * time.Millisecond)
		total++
	}
	println("send total", total)
}
