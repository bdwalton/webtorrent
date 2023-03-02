import { Component, OnInit } from '@angular/core';

import { TorrentService, Torrent } from '../torrent.service';

import { FileSizeFormatterPipe } from '../file-size-formatter.pipe';

@Component({
  selector: 'app-torrent',
  templateUrl: './torrent.component.html',
  styleUrls: ['./torrent.component.scss']
})
export class TorrentComponent implements OnInit {

  torrents: Torrent[] = [];
  torrentURI: string = '';

  constructor(private torrentService: TorrentService) { }

  ngOnInit() {
    this.getTorrents();
  }

  getTorrents() {
    this.torrentService.getTorrents().subscribe((data: Torrent[]) => {
      this.torrents = data;
    });
  }

  addTorrent() {
    var newTorrent = new Torrent();
    newTorrent.URI = this.torrentURI;

    this.torrentService.addTorrent(newTorrent).subscribe((torrent: Torrent) => {
      this.torrents.push(torrent);
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
    this.torrentService.deleteTorrent(torrent).subscribe((t: Torrent) => {
      const idx = this.torrents.findIndex((elem) => {
        return elem.Hash === t.Hash
      });
      if (idx !== -1) {
        this.torrents.splice(idx, 1)
      }
    })
  }
}
