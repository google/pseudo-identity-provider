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

import {Component} from '@angular/core';

import {ConfigService} from '../config.service';

@Component({
    selector: 'app-edit',
    templateUrl: './edit.component.html',
    styleUrl: './edit.component.css',
    standalone: false
})
export class EditComponent {
  value = '';

  constructor(private configService: ConfigService) {
    this.loadConfig();
  }

  loadConfig() {
    this.configService.getConfig().subscribe({
      next: (resp: any) => (this.value = JSON.stringify(resp, null, 2)),
    });
  }

  onSubmit() {
    this.configService.setConfig(this.value).subscribe();
  }
}
