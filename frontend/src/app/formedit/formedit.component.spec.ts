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
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatDividerModule} from '@angular/material/divider';
import {MatExpansionModule} from '@angular/material/expansion';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {MatListModule} from '@angular/material/list';
import {MatTabsModule} from '@angular/material/tabs';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormlyModule} from '@ngx-formly/core';
import {FormlyMaterialModule} from '@ngx-formly/material';

import {FormeditComponent} from './formedit.component';
import {ArrayTypeComponent} from '../array.type';
import {ConfigService} from '../config.service';
import {ObjectTypeComponent} from '../object.type';
import {FormlyFieldTabs} from '../tabs.type';
import {FakeConfigService} from '../test-utils';

describe('FormEditComponent', () => {
  let component: FormeditComponent;
  let fixture: ComponentFixture<FormeditComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [
        ArrayTypeComponent,
        FormeditComponent,
        ObjectTypeComponent, 
        FormlyFieldTabs],
      imports: [
        BrowserModule,
        BrowserAnimationsModule,
        FormlyMaterialModule,
        FormlyModule.forRoot({
          types: [
            {name: 'array', component: ArrayTypeComponent},
            {name: 'object', component: ObjectTypeComponent},
            {name: 'tabs', component: FormlyFieldTabs},
          ],
        }),
        FormsModule,
        MatButtonModule,
        MatDividerModule,
        MatExpansionModule,
        MatFormFieldModule,
        MatInputModule,
        MatListModule,
        MatTabsModule,
        ReactiveFormsModule],
      providers: [{provide: ConfigService, useClass: FakeConfigService}],
    }).compileComponents();
    
    fixture = TestBed.createComponent(FormeditComponent);
    component = fixture.componentInstance;
  });

  it('loads and displays the config', waitForAsync(() => {
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.model).toEqual({val: 'test'});
    });
  }));

  it('submit updates config', waitForAsync(() => {
    component.model = {val: 'othervalue'}
    fixture.nativeElement.querySelector('#formedit_submit').click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.config).toEqual({val: 'othervalue'});
    });
  }));

  it('resets to default config', waitForAsync(() => {
    component.model = {val: 'othervalue'}
    fixture.nativeElement.querySelector('#formedit_submit').click();
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      expect(component.config).toEqual({val: 'othervalue'});

      fixture.nativeElement.querySelector('#formedit_reset').click();
      fixture.whenStable().then(() => {
        expect(component.config).toEqual({val: 'test'});
      });
    });
  }));
});