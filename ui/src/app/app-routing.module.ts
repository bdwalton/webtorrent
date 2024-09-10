import { RouterModule, Routes } from '@angular/router';
import { NgModule } from '@angular/core';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { TorrentComponent } from './torrent/torrent.component';
import { TorrentDetailsComponent } from './torrent-details/torrent-details.component';
import { TorrentStatusComponent } from './torrentstatus/torrentstatus.component';

const routes: Routes = [
  {
    path: '',
    redirectTo: 'torrent',
    pathMatch: 'full',
  },
  {
    path: 'torrent',
    component: TorrentComponent,
    title: 'WebTorrent - Torrents',
  },
  {
    path: 'torrentdetails/:id',
    component: TorrentDetailsComponent,
    title: 'WebTorrent - Torrent Details',
  },
  {
    path: 'torrentstatus',
    component: TorrentStatusComponent,
    title: 'WebTorrent - Torrent Client Status',
  },
  {
    path: '**',
    pathMatch: 'full',
    component: PageNotFoundComponent,
    title: 'Page Not Found',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
