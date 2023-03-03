import { Component, OnInit } from '@angular/core';
import { TorrentService, TorrentTextData } from '../torrent.service';

@Component({
  selector: 'app-torrentstatus',
  templateUrl: './serverconfig.component.html',
  styleUrls: ['./serverconfig.component.scss']
})
export class ServerConfigComponent implements OnInit {

  config: string = '';

  constructor(private torrentService: TorrentService) { }

  ngOnInit() {
    this.getStatus();
  }

  getStatus() {
    this.torrentService.getConfig().subscribe((data: TorrentTextData) => {
      this.config = data.Data
    })
  }
}
