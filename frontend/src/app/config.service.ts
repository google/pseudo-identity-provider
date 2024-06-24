/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import {HttpClient, HttpErrorResponse, HttpHeaders} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Observable, throwError} from 'rxjs';
import {catchError, retry, tap} from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {
  configEndpoint = '/config';
  configSchemaEndpoint = '/configschema';
  schema: any = '';

  constructor(private http: HttpClient) {}

  getConfig() {
    return this.http.get<any>(this.configEndpoint).pipe(
      retry(3),
      catchError(this.handleError),
    );
  }

  getConfigSchema() {
    const pipe = this.http.get<any>(this.configSchemaEndpoint).pipe(
      retry(3),
      catchError(this.handleError),
      tap((schema) => this.schema = schema),
    );
    return pipe;
  }

  getCachedSchema() {
    return this.schema;
  }

  setConfig(config: string) {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
        // CSRF defense with Custom Request headers.
        'X-Pseudo-IDP-CSRF-Protection': '1',
      }),
    };

    return this.http
      .post(this.configEndpoint, config, httpOptions)
      .pipe(catchError(this.handleError));
  }

  resetConfig() {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
        // CSRF defense with Custom Request headers.
        'X-Pseudo-IDP-CSRF-Protection': '1',
      }),
    };

    return this.http
      .delete(this.configEndpoint, httpOptions)
      .pipe(catchError(this.handleError));
  }

  private handleError(error: HttpErrorResponse) {
    if (error.status === 0) {
      console.error('An error occurred:', error.error);
    } else {
      console.error(
        `Backend returned code ${error.status}, body was: `,
        error.error,
      );
    }
    return throwError(() => new Error('Unable to query backend.'));
  }
}
