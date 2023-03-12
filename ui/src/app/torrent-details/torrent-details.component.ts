import { LiveAnnouncer } from '@angular/cdk/a11y';
import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { MatSort, Sort } from '@angular/material/sort';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { MatTableDataSource } from '@angular/material/table';
import {
  TorrentService,
  TorrentDetails,
  TorrentFile,
} from '../torrent.service';
import {
  Action,
  TorrentAction,
} from '../torrent-controls/torrent-controls.component';

@Component({
  selector: 'app-torrent-details',
  templateUrl: './torrent-details.component.html',
  styleUrls: ['./torrent-details.component.scss'],
})
export class TorrentDetailsComponent implements OnInit, AfterViewInit {
  id: string = ''; // The unique id for the torrent
  interval: number = 0; // The refresh interval (0 is disabled)
  torrent: TorrentDetails = new TorrentDetails();
  torrentFiles = new MatTableDataSource<TorrentFile>([]);

  displayedColumns: string[] = ['Position', 'Path'];

  constructor(
    private _route: ActivatedRoute,
    private _router: Router,
    private _sanitizer: DomSanitizer,
    private _torrentService: TorrentService,
    private _liveAnnouncer: LiveAnnouncer
  ) {}

  // Must be set in the component html or this will be undefined.
  @ViewChild(MatSort) sort!: MatSort;

  ngOnInit() {
    this._route.paramMap.subscribe((params: ParamMap) => {
      this.id = params.get('id') as string;
    });

    this.getTorrentDetails(this.id);

    this.interval = setInterval(() => {
      this.getTorrentDetails(this.id);
    }, 5000);
  }

  ngOnDestroy() {
    if (this.interval) {
      clearInterval(this.interval);
    }
  }

  ngAfterViewInit() {
    this.torrentFiles.sort = this.sort;
  }

  sanitize(url: string) {
    return this._sanitizer.bypassSecurityTrustUrl(url);
  }

  /** Announce the change in sort state for assistive technology. */
  announceSortChange(sortState: Sort) {
    if (sortState.direction) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  getTorrentDetails(id: string) {
    this._torrentService
      .getTorrentDetails(id)
      .subscribe((data: TorrentDetails) => {
        this.torrent = data;
        this.torrentFiles.data = [...this.torrent.Files];
      });
  }

  handleControlAction(ta: TorrentAction) {
    if (ta.action == Action.DELETE) {
      this._router.navigate(['/torrent']);
    }
  }
}
