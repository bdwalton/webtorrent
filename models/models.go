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

type TorrentFile struct {
	Path string `json:path`
}

type TorrentData struct {
	BasicTorrentData
	Magnet string        `json:magnet`
	Error  string        `json:error`
	Files  []TorrentFile `json:files`
}

func buildTorrentFiles(t *torrent.Torrent) []TorrentFile {
	fp, err := t.FilePaths()
	// When the torrent metadata hasn't been retrieved yet.
	if err != nil {
		return []TorrentFile{}
	}

	tf := make([]TorrentFile, 0, len(fp))
	for _, tfp := range fp {
		tf = append(tf, TorrentFile{Path: tfp})
	}

	return tf
}

func TorrentDataFromTorrent(t *torrent.Torrent) TorrentData {
	s := t.Stats()
	e := ""
	if s.Error != nil {
		e = s.ETA.String()
	}

	// Private torrents will return an error here. We will just
	// pass the empty string and let the client handle that case.
	m, _ := t.Magnet()

	return TorrentData{
		BasicTorrentData: BasicTorrentDataFromTorrent(t),
		Magnet:           m,
		Error:            e,
		Files:            buildTorrentFiles(t),
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
