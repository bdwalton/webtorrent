import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { TorrentService } from './torrent.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  title = 'WebTorrent';

  constructor(
    private _torrentService: TorrentService,
    private _router: Router,
  ) {}

  onClick(event: Event): void {
    event.preventDefault();
    this._torrentService.signOut().subscribe((any) => {});
    this._router.navigate(['/']);
  }
}
