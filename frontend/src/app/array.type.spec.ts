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

import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import {TestBed} from '@angular/core/testing';
import {Component} from '@angular/core';
import {MatExpansionModule} from '@angular/material/expansion';
import {MatInputModule} from '@angular/material/input';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormGroup, ReactiveFormsModule, FormArray} from '@angular/forms';
import {FormlyModule, FormlyFieldConfig, FormlyFormOptions} from '@ngx-formly/core';
import {FormlyMaterialModule} from '@ngx-formly/material';
import {ArrayTypeComponent} from './array.type';
import {createGenericTestComponent} from './test-utils';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

@Component({ selector: 'formly-form-test', template: ''})
class TestComponent {
  fields = app.fields;
  form = app.form;
  model = app.model;
  options = app.options;
}

let app: Partial<{
  form: FormGroup | FormArray;
  fields: FormlyFieldConfig[];
  options: FormlyFormOptions;
  model: any;
}>;

describe('Array component', () => {
  beforeEach(() => {
    app = {
      form: new FormGroup({}),
      model: {},
    };
    TestBed.configureTestingModule({
    declarations: [TestComponent, ArrayTypeComponent],
    imports: [ReactiveFormsModule,
        FormlyMaterialModule,
        MatExpansionModule,
        MatInputModule,
        BrowserAnimationsModule,
        FormlyModule.forRoot({
            types: [
                {
                    name: 'array',
                    component: ArrayTypeComponent,
                },
            ],
        })],
    providers: [provideHttpClient(withInterceptorsFromDi()), provideHttpClientTesting()]
});
  });

  it('should show content', () => {
    app.model = {array: [{value: 'value1'}, {value: 'value2'}]};
    app.fields = [{
      key: 'array',
      type: 'array',
      props: { label: 'title', description: 'arrDesc' },
      fieldArray: { fieldGroup: [{key: 'value', type: 'input'}] },
    }];

    const fixture = createGenericTestComponent('<formly-form [form]="form" [fields]="fields" [model]="model" [options]="options"></formly-form>', TestComponent);
    fixture.detectChanges();
    const title = fixture.nativeElement.querySelector('mat-panel-title');
    expect(title).toBeTruthy();
    expect(title.innerText).toEqual('title');

    const description = fixture.nativeElement.querySelector('mat-panel-description');
    expect(description).toBeTruthy();
    expect(description.innerText).toEqual('arrDesc');

    const entries = fixture.nativeElement.querySelectorAll('.form-array-field');
    expect(entries.length).toEqual(2);
  });

  it('should add an entry', () => {
    app.model = {array: [{value: 'value1'}, {value: 'value2'}]};
    app.fields = [{
      key: 'array',
      type: 'array',
      props: { label: 'title', description: 'arrDesc' },
      fieldArray: { fieldGroup: [{key: 'value', type: 'input'}] },
    }];

    const fixture = createGenericTestComponent('<formly-form [form]="form" [fields]="fields" [model]="model" [options]="options"></formly-form>', TestComponent);
    fixture.nativeElement.querySelector('#array_add_entry').click();
    fixture.detectChanges();
    const entries = fixture.nativeElement.querySelectorAll('.form-array-field');
    expect(entries.length).toEqual(3);
  });

  it('should remove an entry', () => {
    app.model = {array: [{value: 'value1'}, {value: 'value2'}]};
    app.fields = [{
      key: 'array',
      type: 'array',
      props: { label: 'title', description: 'arrDesc' },
      fieldArray: { fieldGroup: [{key: 'value', type: 'input'}] },
    }];

    const fixture = createGenericTestComponent('<formly-form [form]="form" [fields]="fields" [model]="model" [options]="options"></formly-form>', TestComponent);
    fixture.nativeElement.querySelector('#array_remove_entry').click();
    fixture.detectChanges();
    const entries = fixture.nativeElement.querySelectorAll('.form-array-field');
    expect(entries.length).toEqual(1);
  });
});