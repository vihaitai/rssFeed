package slack

import "testing"

func TestSend(t *testing.T) {
	payload := `
	{
		"type": "mrkdwn",
		"text": "<https://zhuanlan.zhihu.com/p/91383212|Wireguard：简约之美>"
	}
	`
	Notify(payload)
}
