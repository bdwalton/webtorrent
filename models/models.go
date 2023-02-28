package models

type Torrent struct {
	URI  string `json:uri`
	Name string `json:name`
	Hash string `json:hash`
}
