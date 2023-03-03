import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { TorrentComponent } from './torrent/torrent.component';
import { TorrentStatusComponent } from './torrentstatus/torrentstatus.component';

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
    path: 'torrentstatus',
    component: TorrentStatusComponent,
    title: "WebTorrent - Torrent Client Status"
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
