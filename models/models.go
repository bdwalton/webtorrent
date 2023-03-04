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

type TextData struct {
	Data string `json:data`
}

func TextDataFromString(d string) *TextData {
	return &TextData{Data: d}
}

// BasicMetaData is an internal datatype for maintaining state that
// the torrent library doesn't or shouldn't.
type BasicMetaData struct {
	// The torrent library doesn't maintain this after we consume
	// it to start a new torrent. We think it's useful, so we'll
	// hang onto it.
	URI string
	// The torrent library doesn't have a notion of pause or
	// ability to query whether it is running or not. To do that,
	// need to maintain our own view of whe torrent.
	Running bool
	// A reference to the torrent itself
	T *torrent.Torrent
}
