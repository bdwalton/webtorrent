import { LiveAnnouncer } from '@angular/cdk/a11y';
import {
  AfterViewInit,
  HostListener,
  Component,
  OnInit,
  ViewChild,
} from '@angular/core';
import { MatSort, Sort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { MatSnackBar } from '@angular/material/snack-bar';
import { filter } from 'rxjs/operators';
import { TorrentService, Torrent } from '../torrent.service';
import { FileSizeFormatterPipe } from '../file-size-formatter.pipe';
import {
  MatDialog,
  MAT_DIALOG_DATA,
  MatDialogRef,
} from '@angular/material/dialog';
import {
  AddTorrentDialogComponent,
  DialogData,
} from '../add-torrent-dialog/add-torrent-dialog.component';
import {
  Action,
  TorrentAction,
} from '../torrent-controls/torrent-controls.component';

@Component({
  selector: 'app-torrent',
  templateUrl: './torrent.component.html',
  styleUrls: ['./torrent.component.scss'],
})
export class TorrentComponent implements OnInit, AfterViewInit {
  autoRefresh: boolean = true;
  interval: number = 0;
  torrents = new MatTableDataSource<Torrent>([]);

  displayedColumns: string[] = [
    'Name',
    'NFiles',
    'Status',
    'Progress',
    'Controls',
  ];

  constructor(
    public dialog: MatDialog,
    private _snackbar: MatSnackBar,
    private _torrentService: TorrentService,
    private _liveAnnouncer: LiveAnnouncer
  ) {}

  // Must be set in the component html or this will be undefined.
  @ViewChild(MatSort) sort!: MatSort;

  ngOnInit() {
    this.getTorrents();
    this.interval = setInterval(() => {
      this.getTorrents();
    }, 5000);
  }

  ngAfterViewInit() {
    this.torrents.sort = this.sort;
    this.torrents.sortingDataAccessor = (
      row: Torrent,
      columnName: string
    ): string | number => {
      if (columnName == 'Progress') {
        return (
          (row.TotalProgress.BytesDown / row.TotalProgress.BytesTotal) * 100
        );
      }

      return row[columnName as keyof Torrent].toString();
    };
  }

  /** Announce the change in sort state for assistive technology. */
  announceSortChange(sortState: Sort) {
    if (sortState.direction) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  @HostListener('window:keydown.a', ['$event'])
  triggerDialog(event: KeyboardEvent) {
    event.preventDefault();
    this.addTorrentDialog();
  }

  toggleRefresh(): void {
    this.autoRefresh = !this.autoRefresh;
    if (this.autoRefresh) {
      this.interval = setInterval(() => {
        this.getTorrents();
      }, 5000);
      this._snackbar.open('Refresh enabled every 5s.', 'Ok', {
        duration: 5000,
      });
    } else {
      clearInterval(this.interval);
      this._snackbar.open('Refresh disabled.', 'Ok', { duration: 5000 });
    }
  }

  addTorrentDialog(): void {
    const dialogRef = this.dialog.open(AddTorrentDialogComponent, {
      data: { uri: '' },
    });

    dialogRef.afterClosed().subscribe((result) => {
      // If we cancel the dialog, it's closed and this is still
      // called. Ensure we have a useable result before we inspect it.
      if (typeof result === 'undefined') {
        return;
      }

      this.addTorrent(result);
    });
  }

  getTorrents() {
    this._torrentService.getTorrents().subscribe((data: Torrent[]) => {
      this.torrents.data = [...data];
    });
  }

  addTorrent(uri: string) {
    this._torrentService.addTorrent(uri).subscribe((data: Torrent) => {
      // We can add a torrent with an identical hash, but the backend
      // doesn't consider that an error and will return it
      // happily. Thus, we ensure that we don't add dups to our list.
      var idx = this.torrents.data.findIndex((item) => item.Hash === data.Hash);
      if (idx === -1) {
        this.torrents.data = [...this.torrents.data, data];
      }
    });
  }

  handleControlAction(ta: TorrentAction) {
    if (ta.action == Action.DELETE) {
      this.torrents.data = this.torrents.data.filter(
        (item) => item.ID != ta.torrent.ID
      );
    } else {
      // Filter the old torrent and re-add with updated data from the
      // server.
      var newdata = this.torrents.data.filter(
        (item) => item.ID != ta.torrent.ID
      );
      this.torrents.data = [...newdata, ta.torrent];
    }
  }
}
