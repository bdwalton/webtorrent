import { Component, OnInit } from '@angular/core';

import { TorrentService, TorrentTextData } from '../torrent.service';

@Component({
  selector: 'app-torrentstatus',
  templateUrl: './torrentstatus.component.html',
  styleUrls: ['./torrentstatus.component.scss']
})
export class TorrentStatusComponent implements OnInit {

  status: string = '';

  constructor(private torrentService: TorrentService) { }

  ngOnInit() {
    this.getStatus();
  }

  getStatus() {
    this.torrentService.getStatus().subscribe((data: TorrentTextData) => {
      this.status = data.Data
    })
  }
}