package models

import "github.com/cenkalti/rain/torrent"

type Progress struct {
	BytesDown  int64 `json:bytesdown`
	BytesTotal int64 `json:bytestotal`
}

type BasicTorrentData struct {
	ID            string   `json:id`
	Hash          string   `json:hash`
	Name          string   `json:name`
	Status        string   `json:status`
	TotalProgress Progress `json:progress`
}

func BasicTorrentDataFromTorrent(t *torrent.Torrent) *BasicTorrentData {
	s := t.Stats()
	return &BasicTorrentData{
		ID:            t.ID(),
		Hash:          t.InfoHash().String(),
		Name:          t.Name(),
		Status:        s.Status.String(),
		TotalProgress: Progress{s.Bytes.Completed, s.Bytes.Total},
	}
}

type TorrentURI struct {
	URI string `json:uri`
}

type TorrentID struct {
	ID string `json:id`
}

type APIError struct {
	Details string `json:details`
}

type ServerData struct {
	Data string `json:data`
}
