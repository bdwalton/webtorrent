import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class TorrentService {
  constructor(private httpClient: HttpClient) {}

  getTorrents() {
    return this.httpClient.get<Torrent[]>(environment.gateway + '/v1/torrent');
  }

  addTorrent(uri: string) {
    var turi = new TorrentURI(uri);
    return this.httpClient.post<Torrent>(
      environment.gateway + '/v1/torrent',
      turi
    );
  }

  startTorrent(id: string) {
    var tid = new TorrentID(id);
    return this.httpClient.put<Torrent>(
      environment.gateway + '/v1/torrent/start',
      tid
    );
  }

  pauseTorrent(id: string) {
    var tid = new TorrentID(id);
    return this.httpClient.put<Torrent>(
      environment.gateway + '/v1/torrent/pause',
      tid
    );
  }

  deleteTorrent(id: string) {
    return this.httpClient.delete<Torrent>(
      environment.gateway + '/v1/torrent/' + id
    );
  }

  getTorrentDetails(id: string) {
    return this.httpClient.get<TorrentDetails>(
      environment.gateway + '/v1/torrentdetails/' + id
    );
  }

  getStatus() {
    return this.httpClient.get<ServerData>(
      environment.gateway + '/v1/torrentstatus'
    );
  }

  getConfig() {
    return this.httpClient.get<ServerData>(
      environment.gateway + '/v1/showconfig'
    );
  }
}

export class Progress {
  BytesDown: number = 0;
  BytesTotal: number = 0;
}

export class TorrentFile {
  Path: string = '';
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
