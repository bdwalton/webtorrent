import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { TorrentComponent } from './torrent/torrent.component';
import { TorrentDetailsComponent } from './torrent-details/torrent-details.component';
import { TorrentStatusComponent } from './torrentstatus/torrentstatus.component';
import { ServerConfigComponent } from './serverconfig/serverconfig.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'torrent',
    pathMatch: 'full'
  },
  {
    path: 'torrent',
    component: TorrentComponent,
    title: "WebTorrent - Torrents"
  },
  {
    path: 'torrentdetails/:hash',
    component: TorrentDetailsComponent,
    title: "WebTorrent - Torrent Details"
  },
  {
    path: 'torrentstatus',
    component: TorrentStatusComponent,
    title: "WebTorrent - Torrent Client Status"
  },
  {
    path: 'showconfig',
    component: ServerConfigComponent,
    title: "WebTorrent - WebTorrent Server Config"
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
