package models

import (
	"github.com/anacrolix/torrent"
)

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

type Torrent struct {
	URI        string `json:uri`
	Name       string `json:name`
	Hash       string `json:hash`
	Running    bool   `json:running`
	Done       bool   `json:done`
	BytesDown  int64  `json:bytesdown`
	BytesTotal int64  `json:bytestotal`
}

func FromTorrentData(md *BasicMetaData) *Torrent {
	t := &Torrent{
		URI:        md.URI,
		Name:       md.T.Name(),
		Hash:       md.T.InfoHash().HexString(),
		Running:    md.Running,
		BytesDown:  md.T.BytesCompleted(),
		BytesTotal: md.T.Length(),
	}
	t.Done = (t.BytesDown == t.BytesTotal)
	return t
}

type TextData struct {
	Data string `json:data`
}

func TextDataFromString(d string) *TextData {
	return &TextData{Data: d}
}

type APIError struct {
	Error  string `json:error`
	Detail string `json:detail`
}
