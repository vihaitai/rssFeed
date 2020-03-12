package slack

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var _webHook string
var _client = http.Client{}

func init() {
	_webHook = os.Getenv("slackHook")
}

func Notify(msg string) error {
	if _webHook == "" {
		log.Printf("[INFO]slackHook env should be provided")
		return nil
	}
	log.Printf("[INFO] slack notify payload %s", msg)
	req, err := http.NewRequest("POST", _webHook, strings.NewReader(msg))
	if err != nil {
		return err
	}
	resp, err := _client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("slack web hook status code %d", resp.StatusCode))
	}
	return nil
}
