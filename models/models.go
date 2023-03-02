package models

import (
	"github.com/anacrolix/torrent"
)

type Torrent struct {
	URI        string `json:uri`
	Name       string `json:name`
	Hash       string `json:hash`
	BytesDown  int64  `json:bytesdown`
	BytesTotal int64  `json:bytestotal`
}

func FromTorrent(t *torrent.Torrent) *Torrent {
	return &Torrent{
		Name:       t.Name(),
		Hash:       t.InfoHash().HexString(),
		BytesDown:  t.BytesCompleted(),
		BytesTotal: t.Length(),
	}
}
