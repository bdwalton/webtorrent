import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { TorrentComponent } from './torrent/torrent.component';

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
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
