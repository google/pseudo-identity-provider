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
import {FormGroup} from '@angular/forms';
import {FormlyFieldConfig} from '@ngx-formly/core';
import {FormlyJsonschema} from '@ngx-formly/core/json-schema';
import {JSONSchema7} from 'json-schema';

import {ConfigService} from '../config.service';
import * as helpers from '../formly.helpers';

@Component({
    selector: 'app-formedit',
    templateUrl: './formedit.component.html',
    styleUrl: './formedit.component.css',
    standalone: false
})
export class FormeditComponent {
  form: FormGroup = new FormGroup({});
  model: any = {};
  fields: FormlyFieldConfig[] = [];
  schema: any = '';
  config: any = '';

  constructor(
    private formlyJsonschema: FormlyJsonschema,
    private configService: ConfigService,
  ) {
    this.loadExample();
  }

  loadExample() {
    this.configService.getConfigSchema().subscribe({
      next: (resp: any) => (this.schema = resp),
      complete: () => this.loadConfig(),
    });
  }

  loadConfig() {
    this.configService.getConfig().subscribe({
      next: (resp: any) => (this.config = resp),
      complete: () => this.createForm(),
    });
  }

  createForm() {
    this.form = new FormGroup({});
    const fieldConfig = this.formlyJsonschema.toFieldConfig(this.schema);
    helpers.findAndInsertExpressions(fieldConfig, this.schema);
    this.fields = [
      {
        type: 'tabs',
        fieldGroup: fieldConfig.fieldGroup,
      },
    ];
    this.model = this.config;
  }

  onSubmit() {
    if (this.form.valid) {
      this.configService
        .setConfig(JSON.stringify(this.model, null, 2))
        .subscribe((resp: any) => (this.config = resp));
    }
  }

  onReset() {
    this.configService.resetConfig().subscribe({
      next: (resp: any) => (this.config = resp),
      complete: () => this.createForm(),
    });
  }
}
