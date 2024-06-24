// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package idp

import (
	"bytes"
	"customidp/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	testConfig = &config.Config{
		AuthAction: config.AuthAction{
			Action: "redirect",
			Redirect: config.AuthRedirect{
				DefaultParamAction: "passthrough",
				Parameters: []config.Parameter{
					{ID: "access_token", Action: "random"},
					{ID: "redirect_uri", Action: "omit"},
					{ID: "token_type", Action: "set", Values: []string{"Bearer"}},
				},
				UseHashFragment: true,
			},
		},
	}
)

func TestConfigSchemaHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(configSchemaHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("configSchemaHandler() returned unexpected failure code %d", rr.Code)
	}
}

func TestConfigHandlerGet(t *testing.T) {
	config.SetGlobalConfig(testConfig)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(configHandler)

	handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("configHandler() returned unexpected failure code %d", rr.Code)
	}

	gotConfig := &config.Config{}
	if err := json.NewDecoder(rr.Body).Decode(gotConfig); err != nil {
		t.Fatalf("configHandler() returned unparsable JSON %v", err)
	}

	if !reflect.DeepEqual(testConfig, gotConfig) {
		t.Fatalf("configHandler() expected %v, got %v", testConfig, gotConfig)
	}
}

func TestConfigHandlerPost(t *testing.T) {
	setupCreds(t)

	testConfig := &config.Config{
		AuthAction: config.AuthAction{
			Action: "redirect",
			Redirect: config.AuthRedirect{
				DefaultParamAction: "passthrough",
				Parameters: []config.Parameter{
					{ID: "access_token", Action: "random"},
					{ID: "redirect_uri", Action: "omit"},
					{ID: "token_type", Action: "set", Values: []string{"Bearer"}},
				},
				UseHashFragment: true,
			},
		},
	}

	cases := []struct {
		title    string
		username string
		password string
		wantCode int
		headers  map[string]string
		config   *config.Config
		method   string
	}{
		{
			title:    "Valid config update",
			username: testDefaultUsername,
			password: testDefaultPassword,
			wantCode: 200,
			headers: map[string]string{
				"X-Pseudo-IDP-CSRF-Protection": "1",
				"Content-Type":                 "application/json",
				"Origin":                       "https://idp.idp",
			},
			config: testConfig,
			method: "POST",
		},
		{
			title:    "Reset Config",
			username: testDefaultUsername,
			password: testDefaultPassword,
			wantCode: 200,
			headers: map[string]string{
				"X-Pseudo-IDP-CSRF-Protection": "1",
				"Content-Type":                 "application/json",
				"Origin":                       "https://idp.idp",
			},
			config: nil,
			method: "DELETE",
		},
		{
			title:    "Missing csrf",
			username: testDefaultUsername,
			password: testDefaultPassword,
			wantCode: 400,
			headers: map[string]string{
				"Content-Type": "application/json",
				"Origin":       "https://idp.idp",
			},
			config: testConfig,
			method: "POST",
		},
		{
			title:    "Wrong content type",
			username: testDefaultUsername,
			password: testDefaultPassword,
			wantCode: 400,
			headers: map[string]string{
				"X-Pseudo-IDP-CSRF-Protection": "1",
				"Content-Type":                 "application/text",
				"Origin":                       "https://idp.idp",
			},
			config: testConfig,
			method: "POST",
		},
		{
			title:    "Wrong origin",
			username: testDefaultUsername,
			password: testDefaultPassword,
			wantCode: 400,
			headers: map[string]string{
				"X-Pseudo-IDP-CSRF-Protection": "1",
				"Content-Type":                 "application/json",
				"Origin":                       "https://idp2.idp",
			},
			config: testConfig,
			method: "POST",
		},
		{
			title:    "Wrong password",
			username: testDefaultUsername,
			password: "wrong",
			wantCode: 401,
			headers: map[string]string{
				"X-Pseudo-IDP-CSRF-Protection": "1",
				"Content-Type":                 "application/json",
				"Origin":                       "https://idp.idp",
			},
			config: testConfig,
			method: "POST",
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			config.SetGlobalConfig(&config.DefaultConfig)
			wantConfig := testConfig
			req := &http.Request{}

			switch tc.method {
			case "POST":
				data, err := json.Marshal(tc.config)
				if err != nil {
					t.Fatal(err)
				}

				req, err = http.NewRequest("POST", "https://idp.idp/config", bytes.NewBuffer(data))
				if err != nil {
					t.Fatal(err)
				}
			case "DELETE":
				var err error
				req, err = http.NewRequest("DELETE", "https://idp.idp/config", nil)
				if err != nil {
					t.Fatal(err)
				}
				wantConfig = &config.DefaultConfig
			}

			req.SetBasicAuth(tc.username, tc.password)
			for key, val := range tc.headers {
				req.Header.Set(key, val)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(configHandler)

			handler.ServeHTTP(rr, req)
			if rr.Code != tc.wantCode {
				t.Fatalf("configHandler() returned unexpected code %d, expected %d", rr.Code, tc.wantCode)
			}

			if rr.Code != 200 {
				return
			}

			gotConfig := config.GetGlobalConfig()
			if !reflect.DeepEqual(wantConfig, gotConfig) {
				t.Errorf("configHandler() expected %v, got %v", testConfig, gotConfig)
			}
		})
	}
}
