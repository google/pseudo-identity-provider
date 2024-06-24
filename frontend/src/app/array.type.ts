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
import {FieldArrayType, FormlyFieldConfig} from '@ngx-formly/core';
import {ConfigService} from './config.service';
import * as helpers from './formly.helpers';

@Component({
  selector: 'formly-array-type',
  templateUrl: './array.type.html',
  styleUrls: ['./array.type.css'],
})
export class ArrayTypeComponent extends FieldArrayType {
  constructor(private configService: ConfigService) {
    super();
  }

  override onPopulate(field: ArrayTypeComponent): void {
    super.onPopulate(field);
    helpers.findAndInsertExpressionsForArray(
      field,
      this.configService.getCachedSchema(),
    );
  }

  addField() {
    // Formly does not seem to set the default values for added array entries,
    // so we'll do it ourselves here.
    let defaultModel = {};
    if (this.field && typeof this.field.fieldArray === 'function') {
      defaultModel = this.setDefaults(this.field.fieldArray(this.field));
    }

    this.add(this.field.fieldGroup?.length, defaultModel);
  }

  setDefaults(field: FormlyFieldConfig): Object | string {
    switch (field.type) {
      case 'array': {
        if (typeof field.fieldArray === 'function') {
          return this.setDefaults(field.fieldArray(field));
        }
        return [];
      }
      case 'object': {
        if (!field.fieldGroup) {
          return {};
        }

        let model: any = {};
        for (const subField of field.fieldGroup) {
          if (typeof subField.key !== 'string') {
            continue;
          }
          model[subField.key] = this.setDefaults(subField);
        }
        return model;
      }
      default: {
        return field.defaultValue;
      }
    }
  }
}
