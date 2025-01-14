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

import * as helpers from './formly.helpers'

import {FormlyFieldConfig, FormlyFieldProps, FieldArrayType} from '@ngx-formly/core';
import {FormlyJsonschema} from '@ngx-formly/core/json-schema';
import {JSONSchema7} from 'json-schema';
import {FormGroup, FormArray} from '@angular/forms';
import {createComponent} from '@ngx-formly/core/testing';
import {Component} from '@angular/core';

// Test Array Component.
@Component({
    selector: 'formly-array-type',
    template: `
    <div *ngFor="let field of field.fieldGroup; let i = index">
      <formly-group [field]="field"></formly-group>
      <button [id]="'remove-' + i" type="button" (click)="remove(i)">Remove</button>
    </div>
    <button id="add" type="button" (click)="add()">Add</button>
  `,
    standalone: false
})
class ArrayTypeComponent extends FieldArrayType {}

// Render form from schema.
const renderComponent = ({ schema, model }: { schema: JSONSchema7; model?: any }) => {
  const field = new FormlyJsonschema().toFieldConfig(schema);

  const options = createComponent<{ field: FormlyFieldConfig }>({
    template: `
      <form [formGroup]="form">
        <formly-form
          [model]="model"
          [fields]="fields"
          [options]="options"
          [form]="form">
        </formly-form>
      </form>
    `,
    inputs: {
      fields: [field],
      model: model || {},
      form: Array.isArray(model) ? new FormArray([]) : new FormGroup({}),
    } as any,
    declarations: [ArrayTypeComponent],
    config: {
      types: [
        { name: 'object', extends: 'formly-group' },
        { name: 'array', component: ArrayTypeComponent },
      ],
    },
  });

  return { ...options, field };
};


// Test JSON Schema.
let schema : any = {
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
    'token_action': {
      'properties': {
        'action_type': {
        },
        // Array Type has items then properties.
        'parameters': {
          'type': 'array',
          'items': {
            'properties': {
              'action': {
              },
              'custom_key': {
                'hide': 'action !== custom',
              }
            }
          }
        },
        'error': {
          'hide': 'action_type !== error',
        }
      }
    },
  },
};

describe('findExpressions', () => {
  it('finds expression map', () => {
    let wantMap : Map<string, Map<string, string>> = new Map([
      ['auth_action.redirect', new Map([['hide', 'action_type !== redirect']])],
      ['auth_action.error', new Map([['hide', 'action_type !== error']])],
      ['token_action.parameters.custom_key', new Map([['hide', 'action !== custom']])],
      ['token_action.error', new Map([['hide', 'action_type !== error']])],
    ]);
    let gotMap = helpers.findExpressions(schema, "");

    expect(gotMap).toEqual(wantMap);
  });
})

describe('findAndInsertExpressions', () => {
  let formlyJsonschema: FormlyJsonschema = new FormlyJsonschema();
  let fieldConfig : FormlyFieldConfig;

  beforeEach(async () => {
    fieldConfig = formlyJsonschema.toFieldConfig(schema);
  });

  it ('finds and inserts expressions on the default form config', () => {
    helpers.findAndInsertExpressions(fieldConfig, schema);

    // Check that expected fields have hide expressions.
    // token_action.parameters.custom_key won't have any expressions sets since it is an empty array.
    expect(fieldConfig?.fieldGroup?.[0]?.fieldGroup?.[1]?.expressions?.['hide']).toBeTruthy(); // auth_action.redirect
    expect(fieldConfig?.fieldGroup?.[0]?.fieldGroup?.[2]?.expressions?.['hide']).toBeTruthy(); // auth_action.error
    expect(fieldConfig?.fieldGroup?.[1]?.fieldGroup?.[2]?.expressions?.['hide']).toBeTruthy(); // token_action.error
  });

  /* Disabled test as ngx-formly's createComponent is broken in Angular 19.
     Fix is in https://github.com/ngx-formly/ngx-formly/pull/4021 but is not released yet.
  it ('finds and inserts expressions on array from configs', () => {
    const {field, fixture} = renderComponent({
      schema: schema,
      model: {
        'token_action': {
          'action_type': 'respond',
          'parameters': [
            {'action': 'custom'},
            {'action': 'set'},
          ],
        },
      },
    });
    let arrayField = field.fieldGroup?.[1].fieldGroup?.[1]!;
    helpers.findAndInsertExpressionsForArray(arrayField, schema);

    for (let i = 0; i < arrayField?.fieldGroup?.length!; i++) {
      // Check that hide expression is set for custom_key field.
      expect(arrayField?.fieldGroup?.[i].fieldGroup?.[1]?.expressions?.['hide']).toBeTruthy();
    }
  });*/
})