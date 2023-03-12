import { Component, Input, Output, EventEmitter } from '@angular/core';
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
  @Output() torrentAction = new EventEmitter<TorrentAction>();

  constructor(
    public dialog: MatDialog,
    private _torrentService: TorrentService
  ) {}

  startTorrent(id: string) {
    this._torrentService.startTorrent(id).subscribe((torrent: Torrent) => {
      this.torrentAction.emit(new TorrentAction(torrent, Action.STARTED));
    });
  }

  pauseTorrent(id: string) {
    this._torrentService.pauseTorrent(id).subscribe((torrent: Torrent) => {
      this.torrentAction.emit(new TorrentAction(torrent, Action.PAUSED));
    });
  }

  deleteTorrent(id: string) {
    this._torrentService.deleteTorrent(id).subscribe((torrent: Torrent) => {
      this.torrentAction.emit(new TorrentAction(torrent, Action.DELETE));
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

export const enum Action {
  NOOP,
  DELETE,
  STARTED,
  PAUSED,
}

export class TorrentAction {
  torrent: Torrent = new Torrent();
  action: Action = Action.NOOP;

  public constructor(torrent: Torrent, action: Action) {
    this.torrent = torrent;
    this.action = action;
  }
}
