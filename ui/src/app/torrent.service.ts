import { Injectable } from '@angular/core';

import { HttpClient } from '@angular/common/http';

import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class TorrentService {

  constructor(private httpClient: HttpClient) { }

  getTorrents() {
    return this.httpClient.get<Torrent[]>(environment.gateway + '/v1/torrent')
  }

  addTorrent(uri: string) {
    var uriData = new TorrentTextData(uri);
    return this.httpClient.post<Torrent>(environment.gateway + '/v1/torrent', uriData);
  }

  startTorrent(hash: string) {
    var hashData = new TorrentTextData(hash);
    return this.httpClient.put(environment.gateway + '/v1/torrent/start', hashData);
  }

  pauseTorrent(hash: string) {
    var hashData = new TorrentTextData(hash);
    return this.httpClient.put(environment.gateway + '/v1/torrent/pause', hashData);
 }

  deleteTorrent(hash: string) {
    return this.httpClient.delete<Torrent>(environment.gateway + '/v1/torrent/' + hash);
  }

  getTorrentDetails(hash: string) {
    return this.httpClient.get<TorrentDetails>(environment.gateway + '/v1/torrentdetails/' + hash)
  }

  getStatus() {
    return this.httpClient.get<TorrentTextData>(environment.gateway + '/v1/torrentstatus');
  }

  getConfig() {
    return this.httpClient.get<TorrentTextData>(environment.gateway + '/v1/showconfig');
  }
}

export class Torrent {
  URI: string = '';
  Hash: string = '';
  Name: string = '';
  Running: boolean = false;
  Done: boolean = false;
  BytesDown: number = 0;
  BytesTotal: number = 0;
}

export class TorrentTextData {
  public constructor(uri: string) {
    this.Data = uri;
  }
  Data: string = '';
}

export class TorrentFile {
  Path: string = '';
  BytesDown: number = 0;
  BytesTotal: number = 0;
}

export class TorrentDetails extends Torrent {
  public constructor() {
    super();
  }

  Files: TorrentFile[] = [];
}
