import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse,
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Injectable()
export class HttpErrorsInterceptor implements HttpInterceptor {
  constructor(private _snackBar: MatSnackBar) {}

  intercept(
    request: HttpRequest<unknown>,
    next: HttpHandler
  ): Observable<HttpEvent<unknown>> {
    return next.handle(request).pipe(
      catchError((err) => {
        // Note that we don't do logging here. We leave that to
        // default event handling. We only intercept errors so we can
        // show a snack bar to the user.
        if (err.error instanceof ErrorEvent) {
          // Client side.
          this._snackBar.open('Local error: ' + err.error, 'Ok', {
            duration: 5000,
          });
        } else {
          // Server side.
          this._snackBar.open('Remote error: ' + err.error.Detail, 'Ok', {
            duration: 5000,
          });
        }

        // We only wanted to do user notification cleanly, so rethrow
        // the error for default, upstream handling.
        return throwError(err);
      })
    );
  }
}
