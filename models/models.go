package models

type TorrentAddress struct {
	URI string `json:"uri"`
}

type TorrentID struct {
	Hash string `json:hash`
}

type TorrentInfo struct {
	Name string `json:name`
	Hash string `json:hash`
}
