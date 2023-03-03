import { Component, OnInit } from '@angular/core';

import { filter } from 'rxjs/operators';

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

  displayedColumns: string[] = ['name', 'progress', 'controls'];

  constructor(private torrentService: TorrentService) { }

  ngOnInit() {
    this.getTorrents();
  }

  getTorrents() {
    this.torrentService.getTorrents().subscribe((data: Torrent[]) => {
      this.torrents = [...data];
    });
  }

  addTorrent() {
    var newTorrent = new Torrent();
    newTorrent.URI = this.torrentURI;

    this.torrentService.addTorrent(newTorrent).subscribe((data: Torrent) => {
      // We can add a torrent with an identical hash, but the backend
      // doesn't consider that an error and will return it
      // happily. Thus, we ensure that we don't add dups to our list.
      var idx = this.torrents.findIndex(item => item.Hash === data.Hash)
      if (idx === -1) {
          this.torrents = [...this.torrents, data];
      }
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
    this.torrentService.deleteTorrent(torrent).subscribe((data: Torrent) => {
      this.torrents = this.torrents.filter(item => item.Hash !== data.Hash);
    })
  }
}
