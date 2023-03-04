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

  addTorrent(data: TorrentTextData) {
    return this.httpClient.post<Torrent>(environment.gateway + '/v1/torrent', data);
  }

  startTorrent(hash: TorrentTextData) {
    return this.httpClient.put(environment.gateway + '/v1/torrent/start', hash);
  }

  pauseTorrent(hash: TorrentTextData) {
    return this.httpClient.put(environment.gateway + '/v1/torrent/pause', hash);
 }

  deleteTorrent(torrent: Torrent) {
    return this.httpClient.delete<Torrent>(environment.gateway + '/v1/torrent/' + torrent.Hash);
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
