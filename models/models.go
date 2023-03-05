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
	URI           string    `json:uri`
	Name          string    `json:name`
	Hash          string    `json:hash`
	Running       bool      `json:running`
	Done          bool      `json:done`
	TotalProgress *Progress `json:bytestotal`
}

func FromTorrentData(md *BasicMetaData) *Torrent {
	t := &Torrent{
		URI:           md.URI,
		Name:          md.T.Name(),
		Hash:          md.T.InfoHash().HexString(),
		Running:       md.Running,
		TotalProgress: newProgress(md.T.BytesCompleted(), md.T.Length()),
	}
	t.Done = (t.TotalProgress.BytesDown == t.TotalProgress.BytesTotal)
	return t
}

// torrentFile is an internal container for storing additional info
// about each file in a torrent.
type torrentFile struct {
	Path         string    `json:path`
	FileProgress *Progress `json:fileprogress`
}

// TorrentDetails is what we'll serialize to json to service requests
// for more specific info about a torrent.
type TorrentDetails struct {
	Torrent
	Files []*torrentFile `json:files`
}

type Progress struct {
	BytesDown  int64 `json:bytesdown`
	BytesTotal int64 `json:bytestotal`
}

func newProgress(d, t int64) *Progress {
	return &Progress{BytesDown: d, BytesTotal: t}
}

func TorrentDetailsFromData(md *BasicMetaData) *TorrentDetails {
	td := &TorrentDetails{
		Torrent: *FromTorrentData(md),
		Files:   make([]*torrentFile, 0),
	}

	for _, f := range md.T.Files() {
		fd := &torrentFile{
			Path:         f.Path(),
			FileProgress: newProgress(f.BytesCompleted(), f.Length()),
		}
		td.Files = append(td.Files, fd)
	}

	return td
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
