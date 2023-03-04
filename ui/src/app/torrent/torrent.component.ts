import { LiveAnnouncer } from '@angular/cdk/a11y';
import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { MatSort, Sort } from '@angular/material/sort';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatTableDataSource } from '@angular/material/table';
import { filter } from 'rxjs/operators';
import { TorrentService, Torrent, TorrentTextData } from '../torrent.service';
import { FileSizeFormatterPipe } from '../file-size-formatter.pipe';
import { MatDialog, MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { AddTorrentDialogComponent, DialogData } from '../add-torrent-dialog/add-torrent-dialog.component';

@Component({
  selector: 'app-torrent',
  templateUrl: './torrent.component.html',
  styleUrls: ['./torrent.component.scss']
})
export class TorrentComponent implements OnInit, AfterViewInit {

  interval: number = 0;
  torrents = new MatTableDataSource<Torrent>([]);

  displayedColumns: string[] = ['Name', 'Progress', 'Controls'];

  constructor(public dialog: MatDialog,
              private torrentService: TorrentService,
              private _snackBar: MatSnackBar,
              private _liveAnnouncer: LiveAnnouncer) { }

  // Must be set in the component html or this will be undefined.
  @ViewChild(MatSort) sort!: MatSort;

  ngOnInit() {
    this.getTorrents();
    this.interval = setInterval(()=>{
      this.getTorrents();
    }, 5000);
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

  openDialog(): void {
    const dialogRef = this.dialog.open(AddTorrentDialogComponent, {
      data: {uri: ''},
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result.startsWith('magnet:')) {
        this.addTorrent(result);
      } else {
        this._snackBar.open('Invalid Magnet URI. Must start with "magnet:"', 'OK', {
          duration: 5000
        });
      }
    });
  }

  getTorrents() {
    this.torrentService.getTorrents().subscribe((data: Torrent[]) => {
      this.torrents.data = [...data];
    });
  }

  addTorrent(uri: string) {
    var newURI = new TorrentTextData(uri);

    this.torrentService.addTorrent(newURI).subscribe((data: Torrent) => {
      // We can add a torrent with an identical hash, but the backend
      // doesn't consider that an error and will return it
      // happily. Thus, we ensure that we don't add dups to our list.
      var idx = this.torrents.data.findIndex(item => item.Hash === data.Hash)
      if (idx === -1) {
          this.torrents.data = [...this.torrents.data, data];
      }
    })
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
