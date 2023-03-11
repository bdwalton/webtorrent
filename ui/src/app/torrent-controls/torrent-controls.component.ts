import { Component, Input } from '@angular/core';
import { Torrent, TorrentService } from '../torrent.service';
import { ConfirmDialogComponent } from '../confirm-dialog/confirm-dialog.component';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'app-torrent-controls',
  templateUrl: './torrent-controls.component.html',
  styleUrls: ['./torrent-controls.component.scss'],
})
export class TorrentControlsComponent {
  @Input() torrent = new Torrent();

  constructor(
    public dialog: MatDialog,
    private _torrentService: TorrentService
  ) {}

  startTorrent(id: string) {
    this._torrentService.startTorrent(id).subscribe(() => {});
  }

  pauseTorrent(id: string) {
    this._torrentService.pauseTorrent(id).subscribe(() => {});
  }

  deleteTorrent(id: string) {
    this._torrentService.deleteTorrent(id).subscribe((data: Torrent) => {
      // this.torrents.data = this.torrents.data.filter(
      //   (item) => item.Hash !== data.Hash
      // );
    });
  }

  confirmDeleteDialog(id: string) {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      disableClose: false,
    });

    dialogRef.componentInstance.confirmMessage =
      'Removing a Torrent deletes the data. Are you sure?';
    dialogRef.afterClosed().subscribe((result: string) => {
      if (result) {
        this.deleteTorrent(id);
      }
    });
  }
}
