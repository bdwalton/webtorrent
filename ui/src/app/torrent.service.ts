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

  addTorrent(torrent: Torrent) {
    return this.httpClient.post<Torrent>(environment.gateway + '/v1/torrent', torrent);
  }

  startTorrent(torrent: Torrent) {
    return this.httpClient.put(environment.gateway + '/v1/torrent/start', torrent);
  }

  pauseTorrent(torrent: Torrent) {
    return this.httpClient.put(environment.gateway + '/v1/torrent/pause', torrent);
  }

  deleteTorrent(torrent: Torrent) {
    return this.httpClient.delete(environment.gateway + '/v1/torrent/' + torrent.Hash);
  }
}

export class Torrent {
  URI: string = '';
  Hash: string = '';
  Name: string = '';
}
