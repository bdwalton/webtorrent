import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { RouteReuseStrategy } from '@angular/router';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { TorrentComponent } from './torrent/torrent.component';
import { TorrentStatusComponent } from './torrentstatus/torrentstatus.component';
import { FileSizeFormatterPipe } from './file-size-formatter.pipe';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { CacheRouteReuseStrategy } from './cache-route-reuse.strategy';
import { AddTorrentDialogComponent } from './add-torrent-dialog/add-torrent-dialog.component';
import { HttpErrorsInterceptor } from './http-errors.interceptor';

import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatCardModule } from '@angular/material/card';
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
import { TorrentDetailsComponent } from './torrent-details/torrent-details.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { ProgressComponentComponent } from './progress-component/progress-component.component';
import { ConfirmDialogComponent } from './confirm-dialog/confirm-dialog.component';
import { TorrentControlsComponent } from './torrent-controls/torrent-controls.component';

@NgModule({
  declarations: [
    AppComponent,
    TorrentComponent,
    TorrentStatusComponent,
    FileSizeFormatterPipe,
    AddTorrentDialogComponent,
    TorrentDetailsComponent,
    PageNotFoundComponent,
    ProgressComponentComponent,
    ConfirmDialogComponent,
    TorrentControlsComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatCardModule,
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
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpErrorsInterceptor,
      multi: true,
    },
    {
      provide: RouteReuseStrategy,
      useClass: CacheRouteReuseStrategy,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
