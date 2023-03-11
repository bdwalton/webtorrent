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
	NumFiles      int      `json:numfiles`
}

func BasicTorrentDataFromTorrent(t *torrent.Torrent) BasicTorrentData {
	s := t.Stats()

	nf := -1
	f, err := t.FilePaths()
	if err == nil {
		nf = len(f)
	}

	return BasicTorrentData{
		ID:            t.ID(),
		Hash:          t.InfoHash().String(),
		Name:          t.Name(),
		Status:        s.Status.String(),
		TotalProgress: Progress{s.Bytes.Completed, s.Bytes.Total},
		NumFiles:      nf,
	}
}

type TorrentData struct {
	BasicTorrentData
	Magnet string   `json:magnet`
	Error  string   `json:error`
	Files  []string `json:files`
}

func TorrentDataFromTorrent(t *torrent.Torrent) TorrentData {
	s := t.Stats()
	e := ""
	if s.Error != nil {
		e = s.ETA.String()
	}

	f, err := t.FilePaths()
	// When the torrent metadata hasn't been retrieved yet.
	if err != nil {
		f = []string{}
	}

	return TorrentData{
		BasicTorrentData: BasicTorrentDataFromTorrent(t),
		Magnet:           t.Magnet(),
		Error:            e,
		Files:            f,
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
