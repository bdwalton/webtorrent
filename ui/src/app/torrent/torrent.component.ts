import { Component, OnInit } from '@angular/core';

import { TorrentService, Torrent } from '../torrent.service';

@Component({
  selector: 'app-torrent',
  templateUrl: './torrent.component.html',
  styleUrls: ['./torrent.component.scss']
})
export class TorrentComponent implements OnInit {

  torrents: Torrent[] = [];

  constructor(private torrentService: TorrentService) { }

  ngOnInit() {
    this.getTorrents();
  }

  getTorrents() {
    this.torrentService.getTorrents().subscribe((data: any) => {
      this.torrents = data;
    });
  }

  // addTorrent() {
  //   var newTorrentAddress : TorrentAddress = {
  //     uri: 'magnet:?xt=urn:btih:KRWPCX3SJUM4IMM4YF5RPHL6ANPYTQPU',
  //   };

  //   this.torrentService.addTorrent(newTorrentAddress).subscribe(() => {
  //     this.getAll();
  //   })
  // }

  deleteTorrent(torrent: Torrent) {
    this.torrentService.deleteTorrent(torrent).subscribe(() => {
      this.getTorrents();
    })
  }
}
