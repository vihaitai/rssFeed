home=/home/ubuntu/rssFeed

server_bin:
	@cd ${home}
	go build -o ${home}/bin/rssServer rssFeed

run_server: server_bin
	sudo ${home}/bin/rssServer &

kill_server:
	ps aux|grep rssServer |grep -v grep|awk '{print $$2}'|xargs sudo kill -9

build_cmd:
	@cd ${home}
	go build -o ${home}/bin/rssSync rssFeed/cmd/sync

.PHONY: build_server build_cmd kill_server run_server
