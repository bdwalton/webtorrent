all: webtorrent

webtorrent: webtorrent-go webtorrent-angular

webtorrent-go: webtorrent-angular
	go build

webtorrent-angular:
	(cd ui; ng build)
