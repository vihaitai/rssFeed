package main

import (
	"os"
	"os/signal"
	"rssFeed/httpServer"
)

func main() {
	// http server来返回rss格式的订阅内容
	go httpServer.Start("0.0.0.0:80")

	// pub server来发布新的文章
	//go pubServer.Start("tcp://0.0.0.0:8001")

	// pub信息源的接收端服务
	//go pubServer.StartREP("tcp://127.0.0.1:8002")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
