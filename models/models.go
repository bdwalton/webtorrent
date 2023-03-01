package models

type Torrent struct {
	URI        string `json:uri`
	Name       string `json:name`
	Hash       string `json:hash`
	BytesDown  int64  `json:bytesdown`
	BytesTotal int64  `json:bytestotal`
}
