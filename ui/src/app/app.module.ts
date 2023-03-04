import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { RouteReuseStrategy } from '@angular/router';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TorrentComponent } from './torrent/torrent.component';
import { TorrentStatusComponent } from './torrentstatus/torrentstatus.component';
import { FileSizeFormatterPipe } from './file-size-formatter.pipe';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CacheRouteReuseStrategy } from './cache-route-reuse.strategy';
import { ServerConfigComponent } from './serverconfig/serverconfig.component';
import { AddTorrentDialogComponent } from './add-torrent-dialog/add-torrent-dialog.component';

import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatToolbarModule } from '@angular/material/toolbar';

@NgModule({
  declarations: [
    AppComponent,
    TorrentComponent,
    TorrentStatusComponent,
    FileSizeFormatterPipe,
    ServerConfigComponent,
    AddTorrentDialogComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatButtonModule,
    MatDialogModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatProgressBarModule,
    MatSidenavModule,
    MatSnackBarModule,
    MatSortModule,
    MatTableModule,
    MatToolbarModule,
  ],
  providers: [{
    provide: RouteReuseStrategy,
    useClass: CacheRouteReuseStrategy,
  }],
  bootstrap: [AppComponent]
})
export class AppModule { }
