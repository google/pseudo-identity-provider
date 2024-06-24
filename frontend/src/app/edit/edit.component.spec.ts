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

import {TestBed, ComponentFixture, waitForAsync} from '@angular/core/testing';
import {FormsModule} from '@angular/forms';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormlyMaterialModule} from '@ngx-formly/material';
import {EditComponent} from './edit.component';
import {ConfigService} from '../config.service';
import {FakeConfigService} from '../test-utils';

describe('EditComponent', () => {
  let component: EditComponent;
  let fixture: ComponentFixture<EditComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [EditComponent],
      imports: [
        BrowserModule,
        BrowserAnimationsModule,
        FormlyMaterialModule,
        FormsModule,
        MatFormFieldModule,
        MatInputModule],
      providers: [{provide: ConfigService, useClass: FakeConfigService}],
    }).compileComponents();
    
    fixture = TestBed.createComponent(EditComponent);
    component = fixture.componentInstance;
  });

  it('loads and displays the config', waitForAsync(() => {
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const text = fixture.nativeElement.querySelector('textarea');
      expect(text).toBeTruthy();
      expect(text.value.replace(/\s/g, "")).toMatch('{"val":"test"}');
    });
  }));

  it('submit updates config', waitForAsync(() => {
    component.value = '{"val": "othervalue"}';
    fixture.nativeElement.querySelector('#edit_form_submit').click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      const text = fixture.nativeElement.querySelector('textarea');
      expect(text).toBeTruthy();
      expect(text.value.replace(/\s/g, "")).toMatch('{"val":"othervalue"}');
    });
  }));
});