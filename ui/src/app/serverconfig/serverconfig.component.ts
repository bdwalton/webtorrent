import { Component, OnInit } from '@angular/core';
import { TorrentService, TorrentTextData } from '../torrent.service';

@Component({
  selector: 'app-torrentstatus',
  templateUrl: './serverconfig.component.html',
  styleUrls: ['./serverconfig.component.scss'],
})
export class ServerConfigComponent implements OnInit {
  config: string = '';

  constructor(private _torrentService: TorrentService) {}

  ngOnInit() {
    this.getStatus();
  }

  getStatus() {
    this._torrentService.getConfig().subscribe((data: TorrentTextData) => {
      this.config = data.Data;
    });
  }
}
