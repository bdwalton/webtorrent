import { Component, OnInit } from '@angular/core';

import { TorrentService, Torrent } from '../torrent.service';

@Component({
  selector: 'app-torrent',
  templateUrl: './torrent.component.html',
  styleUrls: ['./torrent.component.scss']
})
export class TorrentComponent implements OnInit {

  torrents = new Map<String, Torrent>();
  torrentURI: string = '';

  constructor(private torrentService: TorrentService) { }

  ngOnInit() {
    this.getTorrents();
  }

  getTorrents() {
    this.torrentService.getTorrents().subscribe((data: Torrent[]) => {
      var tl = new Map<String, Torrent>();
      data.forEach( (t) => {
        tl.set(t.Hash, t);
      })
      this.torrents = tl;
    });
  }

  addTorrent() {
    var newTorrent = new Torrent();
    newTorrent.URI = this.torrentURI;

    this.torrentService.addTorrent(newTorrent).subscribe((torrent: Torrent) => {
      this.torrents.set(torrent.Hash, torrent);
    })

    this.torrentURI = '';
  }

  startTorrent(torrent: Torrent) {
    this.torrentService.startTorrent(torrent).subscribe(() => {
      this.getTorrents();
    })
  }

  pauseTorrent(torrent: Torrent) {
    this.torrentService.pauseTorrent(torrent).subscribe(() => {
      this.getTorrents();
    })
  }

  deleteTorrent(torrent: Torrent) {
    this.torrentService.deleteTorrent(torrent).subscribe((torrent: Torrent) => {
      this.torrents.delete(torrent.Hash);
    })
  }
}
