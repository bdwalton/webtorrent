import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { TorrentService, TorrentDetails } from '../torrent.service';

@Component({
  selector: 'app-torrent-details',
  templateUrl: './torrent-details.component.html',
  styleUrls: ['./torrent-details.component.scss'],
})
export class TorrentDetailsComponent implements OnInit {
  hash: string = '';
  torrent: TorrentDetails = new TorrentDetails();

  constructor(
    private route: ActivatedRoute,
    private sanitizer: DomSanitizer,
    private torrentService: TorrentService
  ) {}

  ngOnInit() {
    this.route.paramMap.subscribe((params: ParamMap) => {
      this.hash = params.get('hash') as string;
    });

    this.getTorrentDetails(this.hash);
  }

  sanitize(url: string) {
    return this.sanitizer.bypassSecurityTrustUrl(url);
  }

  getTorrentDetails(hash: string) {
    this.torrentService
      .getTorrentDetails(hash)
      .subscribe((data: TorrentDetails) => {
        this.torrent = data;
      });
  }
}
