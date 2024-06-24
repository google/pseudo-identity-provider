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

import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';
import {TestBed} from "@angular/core/testing";
import {ConfigService} from "./config.service";

// Test JSON Schema.
let testSchema : any = {
  '$schema': 'https://json-schema.org/draft/2020-12/schema',
  'properties': {
    'auth_action': {
      'type': 'object',
      'properties': {
        'action_type': {
        },
        'redirect': {
          'hide': 'action_type !== redirect',
        },
        'error': {
          'hide': 'action_type !== error',
        }
      }
    },
  },
};

let testConfig : any = {
  'token_action': {
    'action_type': 'respond',
    'parameters': [
      {'action': 'custom'},
      {'action': 'set'},
    ],
  },
};

let alternateConfig : any = {
  'token_action': {
    'action_type': 'error',
  },
};

describe('ConfigService', () => {
  let configService : ConfigService;
  let controller: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ConfigService],
    });
    configService = TestBed.inject(ConfigService);
    controller = TestBed.inject(HttpTestingController);
  });

  it('fetches the config schema', () => {
    configService.getConfigSchema().subscribe(
      (schema) => {
        expect(schema).toEqual(testSchema);
      }
    );

    const req = controller.expectOne('/configschema');
    req.flush(testSchema);
    controller.verify();
  });

  it('fetches the cached config schema', () => {
    configService.getConfigSchema().subscribe({
      next: (schema) => {
        expect(schema).toEqual(testSchema);
      },
      complete: () => expect(configService.getCachedSchema()).toEqual(testSchema),
    });

    // Expect only one request.
    const req = controller.expectOne('/configschema');
    req.flush(testSchema);
    controller.verify();
  });

  it('fetches the config', () => {
    configService.getConfig().subscribe(
      (config) => {
        expect(config).toEqual(testConfig);
      }
    );

    const req = controller.expectOne('/config');
    req.flush(testConfig);
    controller.verify();
  });

  it('sets the config', () => {
    configService.setConfig(testConfig).subscribe(
      (config) => {
        expect(config).toEqual(testConfig);
      }
    );

    const req = controller.expectOne({method: 'POST', url:'/config'});
    req.flush(testConfig);
    controller.verify();
  });

  it('resets the config', () => {
    configService.resetConfig().subscribe();
    const req = controller.expectOne({method: 'DELETE', url:'/config'});
    expect(req.request.method).toEqual("DELETE");
    req.flush(null);
    controller.verify();
  });
});