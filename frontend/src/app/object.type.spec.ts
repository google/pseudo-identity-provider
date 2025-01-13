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

import {TestBed} from '@angular/core/testing';
import {Component} from '@angular/core';
import {MatExpansionModule} from '@angular/material/expansion';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormGroup, ReactiveFormsModule, FormArray} from '@angular/forms';
import {FormlyModule, FormlyFieldConfig, FormlyFormOptions} from '@ngx-formly/core';
import {ObjectTypeComponent} from './object.type';
import {createGenericTestComponent} from './test-utils';

@Component({
    selector: 'formly-form-test', template: '',
    standalone: false
})
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

describe('Object component', () => {
  beforeEach(() => {
    app = {
      form: new FormGroup({}),
      model: {},
    };
    TestBed.configureTestingModule({
      declarations: [TestComponent, ObjectTypeComponent],
      imports: [
        ReactiveFormsModule,
        MatExpansionModule,
        BrowserAnimationsModule,
        FormlyModule.forRoot({
          types: [
            {
              name: 'object',
              component: ObjectTypeComponent,
            },
          ],
        }),
      ],
    });
  });

  it('should show content', () => {
    app.model = {};
    app.fields = [{
      key: 'object',
      type: 'object',
      props: { label: 'title', description: 'objDesc' },
      fieldGroup: [],
    }];

    const fixture = createGenericTestComponent('<formly-form [form]="form" [fields]="fields" [model]="model" [options]="options"></formly-form>', TestComponent);
    fixture.detectChanges();
    const title = fixture.nativeElement.querySelector('mat-panel-title');
    expect(title).toBeTruthy();
    expect(title.innerText).toEqual('title');

    const description = fixture.nativeElement.querySelector('mat-panel-description');
    expect(description).toBeTruthy();
    expect(description.innerText).toEqual('objDesc');
  });
});