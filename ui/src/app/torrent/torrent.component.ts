import { LiveAnnouncer } from '@angular/cdk/a11y';

import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';

import { MatSort, Sort } from '@angular/material/sort';

import { MatTableDataSource } from '@angular/material/table';

import { filter } from 'rxjs/operators';

import { TorrentService, Torrent } from '../torrent.service';

import { FileSizeFormatterPipe } from '../file-size-formatter.pipe';

@Component({
  selector: 'app-torrent',
  templateUrl: './torrent.component.html',
  styleUrls: ['./torrent.component.scss']
})
export class TorrentComponent implements OnInit, AfterViewInit {

  torrents = new MatTableDataSource<Torrent>([]);

  torrentURI: string = '';

  displayedColumns: string[] = ['Name', 'Progress', 'Controls'];

  constructor(private torrentService: TorrentService,
              private _liveAnnouncer: LiveAnnouncer) { }

  // Must be set in the component html or this will be undefined.
  @ViewChild(MatSort) sort!: MatSort;

  ngOnInit() {
    this.getTorrents();
  }

  ngAfterViewInit() {
    this.torrents.sort = this.sort;
  }

  /** Announce the change in sort state for assistive technology. */
  announceSortChange(sortState: Sort) {
    if (sortState.direction) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  getTorrents() {
    this.torrentService.getTorrents().subscribe((data: Torrent[]) => {
      this.torrents.data = [...data];
    });
  }

  addTorrent() {
    var newTorrent = new Torrent();
    newTorrent.URI = this.torrentURI;

    this.torrentService.addTorrent(newTorrent).subscribe((data: Torrent) => {
      // We can add a torrent with an identical hash, but the backend
      // doesn't consider that an error and will return it
      // happily. Thus, we ensure that we don't add dups to our list.
      var idx = this.torrents.data.findIndex(item => item.Hash === data.Hash)
      if (idx === -1) {
          this.torrents.data = [...this.torrents.data, data];
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
      this.torrents.data = this.torrents.data.filter(item => item.Hash !== data.Hash);
    })
  }
}
