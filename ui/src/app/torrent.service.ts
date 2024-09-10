import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class TorrentService {
  constructor(private httpClient: HttpClient) {}

  signOut() {
    return this.httpClient.get<any>(environment.gateway + 'signout');
  }

  getTorrents() {
    return this.httpClient.get<Torrent[]>(environment.gateway + 'v1/torrent');
  }

  addTorrent(uri: string) {
    var turi = new TorrentURI(uri);
    return this.httpClient.post<Torrent>(
      environment.gateway + 'v1/torrent',
      turi,
    );
  }

  startTorrent(id: string) {
    var tid = new TorrentID(id);
    return this.httpClient.put<Torrent>(
      environment.gateway + 'v1/torrent/start',
      tid,
    );
  }

  pauseTorrent(id: string) {
    var tid = new TorrentID(id);
    return this.httpClient.put<Torrent>(
      environment.gateway + 'v1/torrent/pause',
      tid,
    );
  }

  deleteTorrent(id: string) {
    return this.httpClient.delete<Torrent>(
      environment.gateway + 'v1/torrent/' + id,
    );
  }

  getTorrentDetails(id: string) {
    return this.httpClient.get<TorrentDetails>(
      environment.gateway + 'v1/torrentdetails/' + id,
    );
  }

  getStatus() {
    return this.httpClient.get<ServerStats>(
      environment.gateway + 'v1/torrentstatus',
    );
  }

  getConfig() {
    return this.httpClient.get<ServerData>(
      environment.gateway + 'v1/showconfig',
    );
  }
}

export class Progress {
  BytesDown: number = 0;
  BytesTotal: number = 0;
}

export class TorrentFile {
  Path: string = '';
  FileProgress: Progress = new Progress();
}

export class Torrent {
  ID: string = '';
  URI: string = '';
  Hash: string = '';
  Name: string = '';
  Status: string = '';
  TotalProgress: Progress = new Progress();
  NumFiles: number = -1;
}

export class TorrentDetails extends Torrent {
  public constructor() {
    super();
  }

  Files: TorrentFile[] = [];
  Error: string = '';
  Magnet: string = '';
}

export class TorrentURI {
  public constructor(uri: string) {
    this.URI = uri;
  }
  URI: string = '';
}

export class TorrentID {
  public constructor(id: string) {
    this.ID = id;
  }
  ID: string = '';
}

export class ServerData {
  Data: string = '';
}

export class ServerStats {
  // Time elapsed after creation of the Session object.
  Uptime: number = 0; // A time.Duration in Golang

  // Number of torrents in Session.
  Torrents: number = 0;
  // Total number of connected peers.
  Peers: number = 0;
  // Number of available ports for new torrents.
  PortsAvailable: number = 0;

  // Number of rules in blocklist.
  BlockListRules: number = 0;
  // Time elapsed after the last successful update of blocklist.
  BlockListRecency: number = 0; // A time.Duration in Golang

  // Number of objects in piece read cache.
  // Each object is a block whose size is defined in Config.ReadCacheBlockSize.
  ReadCacheObjects: number = 0;
  // Current size of read cache.
  ReadCacheSize: number = 0;
  // Hit ratio of read cache.
  ReadCacheUtilization: number = 0;

  // Number of reads per second from disk.
  ReadsPerSecond: number = 0;
  // Number of active read requests from disk.
  ReadsActive: number = 0;
  // Number of pending read requests from disk.
  ReadsPending: number = 0;

  // Number of objects in piece write cache.
  // Objects are complete pieces.
  // Piece size differs among torrents.
  WriteCacheObjects: number = 0;
  // Current size of write cache.
  WriteCacheSize: number = 0;
  // Number of pending torrents that is waiting for write cache.
  WriteCachePendingKeys: number = 0;

  // Number of writes per second to disk.
  // Each write is a complete piece.
  WritesPerSecond: number = 0;
  // Number of active write requests to disk.
  WritesActive: number = 0;
  // Number of pending write requests to disk.
  WritesPending: number = 0;

  // Download speed from peers in bytes/s.
  SpeedDownload: number = 0;
  // Upload speed to peers in bytes/s.
  SpeedUpload: number = 0;
  // Read speed from disk in bytes/s.
  SpeedRead: number = 0;
  // Write speed to disk in bytes/s.
  SpeedWrite: number = 0;

  // Number of bytes downloaded from peers.
  BytesDownloaded: number = 0;
  // Number of bytes uploaded to peers.
  BytesUploaded: number = 0;
  // Number of bytes read from disk.
  BytesRead: number = 0;
  // Number of bytes written to disk.
  BytesWritten: number = 0;
}
