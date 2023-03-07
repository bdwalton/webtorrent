import { Component, OnInit } from '@angular/core';

import { TorrentService, ServerData } from '../torrent.service';

@Component({
  selector: 'app-torrentstatus',
  templateUrl: './torrentstatus.component.html',
  styleUrls: ['./torrentstatus.component.scss'],
})
export class TorrentStatusComponent implements OnInit {
  status: string = '';

  constructor(private _torrentService: TorrentService) {}

  ngOnInit() {
    this.getStatus();
  }

  getStatus() {
    this._torrentService.getStatus().subscribe((data: ServerData) => {
      this.status = data.Data;
    });
  }
}
