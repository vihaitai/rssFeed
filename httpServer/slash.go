package httpServer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"rssFeed/db"

	_ "github.com/mattn/go-sqlite3"
	"github.com/slack-go/slack"
)

var signingSecret string

func init() {
	signingSecret = os.Getenv("slackSigningSecret")
}

func handleSlashKeyword(keyword string) ([]byte, error) {
	conn := db.Conn
	query := fmt.Sprintf("select id, title, link, description, created_at, feed_name from articles where title like %%?%% limit 10")
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatalf("prepare %s error %+v", query, err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + keyword + "%")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var msg = slack.Message{}
	for rows.Next() {
		var a = db.Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Link, &a.Description, &a.CreatedAt, &a.FeedName)
		if err != nil {
			log.Printf("article scan error %+v", err)
			continue
		}
		textBlockObject := slack.NewTextBlockObject("mrkdwn", "You have a new test: *Hi there* :wave:", false, false)
		sectionBlock := slack.NewSectionBlock(textBlockObject, nil, nil)
		msg = slack.AddBlockMessage(msg, sectionBlock)
	}
	b, err := json.Marshal(msg)
	return b, err
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
		b, err := handleSlashKeyword(s.Text)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("reponse: %s\n", string(b))
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(s.Command))
		return
	}
}
