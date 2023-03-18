import { Component, OnInit } from '@angular/core';

import { TorrentService, ServerStats } from '../torrent.service';

@Component({
  selector: 'app-torrentstatus',
  templateUrl: './torrentstatus.component.html',
  styleUrls: ['./torrentstatus.component.scss'],
})
export class TorrentStatusComponent implements OnInit {
  _status: ServerStats = new ServerStats();

  _labels: { [item: string]: string[] } = {
    Uptime: ['Uptime', 'Time elapsed after creation of the Session object.'],
    Torrents: ['Torrents', 'Number of torrents in Session.'],
    Peers: ['Peers', 'Total number of connected peers.'],
    PortsAvailable: [
      'Ports Available',
      'Number of available ports for new torrents.',
    ],
    BlockListRules: ['Blocklist Rules', 'Number of rules in blocklist.'],
    BlockListRecency: [
      'Most Recent Blocklist Use',
      'Time elapsed after the last successful update of blocklist.',
    ],
    ReadCacheObjects: [
      'Read Cache Objects',
      'Number of objects in piece read cache. Each object is a block whose size is defined in Config.ReadCacheBlockSize.',
    ],
    ReadCacheSize: ['Read Cache Size', 'Current size of read cache.'],
    ReadCacheUtilization: [
      'Read Cache Utilization',
      'Hit ratio of read cache.',
    ],
    ReadsPerSecond: [
      'Reads Per Second',
      'Number of reads per second from disk.',
    ],
    ReadsActive: ['Reads Active', 'Number of active read requests from disk.'],
    ReadsPending: [
      'Reads Pending',
      'Number of pending read requests from disk.',
    ],
    WriteCacheObjects: [
      'Write Cache Objects',
      'Number of objects in piece write cache. Objects are complete pieces. Piece size differs among torrents.',
    ],
    WriteCacheSize: ['Write Cache Size', 'Current size of write cache.'],
    WriteCachePendingKeys: [
      'Write Cache Pending Keys',
      'Number of pending torrents that are waiting for write cache.',
    ],
    WritesPerSecond: [
      'Writes Per Second',
      'Number of writes per second to disk. Each write is a complete piece.',
    ],
    WritesActive: ['Active Writes', 'Number of active write requests to disk.'],
    WritesPending: [
      'Pending Writes',
      'Number of pending write requests to disk.',
    ],
    SpeedDownload: ['Download Speed', 'Download speed from peers in bytes/s.'],
    SpeedUpload: ['Upload Speed', 'Upload speed to peers in bytes/s.'],
    SpeedRead: ['Read Speed', 'Read speed from disk in bytes/s.'],
    SpeedWrite: ['Write Speed', 'Write speed to disk in bytes/s.'],
    BytesDownloaded: [
      'Bytes Downloaded',
      'Number of bytes downloaded from peers.',
    ],
    BytesUploaded: ['Bytes Uploaded', 'Number of bytes uploaded to peers.'],
    BytesRead: ['Bytes Read', 'Number of bytes read from disk.'],
    BytesWritten: ['Bytes Written', 'Number of bytes written to disk.'],
  };

  displayItems: string[] = [
    'Uptime',
    'Torrents',
    'Peers',
    'PortsAvailable',
    'BlockListRules',
    'BlockListRecency',
    'ReadCacheObjects',
    'ReadCacheSize',
    'ReadCacheUtilization',
    'ReadsPerSecond',
    'ReadsActive',
    'ReadsPending',
    'WriteCacheObjects',
    'WriteCacheSize',
    'WriteCachePendingKeys',
    'WritesPerSecond',
    'WritesActive',
    'WritesPending',
    'SpeedDownload',
    'SpeedUpload',
    'SpeedRead',
    'SpeedWrite',
    'BytesDownloaded',
    'BytesUploaded',
    'BytesRead',
    'BytesWritten',
  ];

  constructor(private _torrentService: TorrentService) {}

  ngOnInit() {
    this.getStatus();
  }

  getStatus() {
    this._torrentService.getStatus().subscribe((data: ServerStats) => {
      this._status = data;
    });
  }

  getTitle(item: string): string {
    return this._labels[item][0];
  }

  getDescription(item: string): string {
    return this._labels[item][1];
  }

  getItem(item: string): number {
    return this._status[item as keyof ServerStats];
  }
}
