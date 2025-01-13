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
import {MatTabsModule} from '@angular/material/tabs';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormGroup, ReactiveFormsModule, FormArray} from '@angular/forms';
import {FormlyModule, FormlyFieldConfig, FormlyFormOptions} from '@ngx-formly/core';
import {FormlyFieldTabs} from './tabs.type';
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

describe('Tabs component', () => {
  beforeEach(() => {
    app = {
      form: new FormGroup({}),
      model: {},
    };
    TestBed.configureTestingModule({
      declarations: [TestComponent, FormlyFieldTabs],
      imports: [
        ReactiveFormsModule,
        MatTabsModule,
        BrowserAnimationsModule,
        FormlyModule.forRoot({
          types: [
            {
              name: 'tab',
              component: FormlyFieldTabs,
            },
          ],
        }),
      ],
    });
  });

  it('should show tab content', () => {
    app.model = {};
    app.fields = [{
      key: 'tab',
      type: 'tab',
      fieldGroup: [
        {
          props: { label: 'tab1' },
          fieldGroup: [
          ],
        },
        {
          props: { label: 'tab2' },
          fieldGroup: [
          ],
        },
      ],
    }];

    const fixture = createGenericTestComponent('<formly-form [form]="form" [fields]="fields" [model]="model" [options]="options"></formly-form>', TestComponent);
    fixture.detectChanges();
    const tabs = fixture.nativeElement.querySelectorAll('.mdc-tab__text-label');
    expect(tabs.length).toEqual(2);
    expect(tabs[0].textContent).toEqual('tab1');
    expect(tabs[1].textContent).toEqual('tab2');
  });
});