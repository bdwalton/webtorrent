import { Component, OnInit } from '@angular/core';

import { TorrentService, Torrent } from '../torrent.service';

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
    var newTorrent : Torrent = {
      URI: this.torrentURI,
      Hash: '',
      Name: '',
    };

    this.torrentService.addTorrent(newTorrent).subscribe((torrent: Torrent) => {
      this.torrents.push(torrent)
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
    this.torrentService.deleteTorrent(torrent).subscribe(() => {
      this.getTorrents();
    })
  }
}
