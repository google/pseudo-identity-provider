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

import {Injectable} from '@angular/core';
import { TestBed, ComponentFixture } from '@angular/core/testing';
import { Observable, of } from 'rxjs';

export function createGenericTestComponent<T>(html: string, type: {new (...args: any[]): T}): ComponentFixture<T> {
  TestBed.overrideComponent(type, {set: {template: html}});
  const fixture = TestBed.createComponent(type);
  fixture.detectChanges();
  return fixture as ComponentFixture<T>;
}

const defaultConfig = '{"val": "test"}';

export
@Injectable()
class FakeConfigService {
  storedConfig = defaultConfig;
  getConfig(): Observable<string> {
    return of(JSON.parse(this.storedConfig));
  }

  setConfig(config: string): Observable<string> {
    this.storedConfig = config;
    return of(JSON.parse(this.storedConfig));
  }

  getConfigSchema() {
    const schema : any = {
      '$schema': 'https://json-schema.org/draft/2020-12/schema',
      'properties': {
        'val': {
          'type': 'string',
        },
      },
    };
    return of(JSON.stringify(schema));
  }

  getCachedSchema() {
    return JSON.parse(this.storedConfig);
  }

  resetConfig() {
    this.storedConfig = defaultConfig;
    return of(JSON.parse(this.storedConfig));
  }
}