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
	"customidp/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUserInfoHandler(t *testing.T) {
	cases := []struct {
		title       string
		config      *config.Config
		wantCode    int
		wantResults map[string]any
	}{
		{
			title:    "Default Config",
			wantCode: 200,
			wantResults: map[string]any{
				"sub":   "12345abcde",
				"email": "testsub@idp.idp",
			},
		},
		{
			title: "Custom Config",
			config: &config.Config{
				UserInfoAction: config.UserInfoAction{
					Action: "respond",
					Respond: config.UserInfoRespond{
						Parameters: []config.Parameter{
							{ID: "custom", Action: "set", Values: []string{"https://{{.Domain}}"}, JSONType: "string"},
						},
					},
				},
			},
			wantCode: 200,
			wantResults: map[string]any{
				"custom": "https://idp.idp",
			},
		},
		{
			title: "Error response",
			config: &config.Config{
				UserInfoAction: config.UserInfoAction{
					Action: "error",
					Error: config.Error{
						ErrorCode:    400,
						ErrorContent: "error",
					},
				},
			},
			wantCode: 400,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			if tc.config != nil {
				config.SetGlobalConfig(tc.config)
			} else {
				config.SetGlobalConfig(&config.DefaultConfig)
			}

			req, err := http.NewRequest("GET", "https://idp.idp/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(userInfoHandler)

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.wantCode {
				t.Fatalf("userInfoHandler() returned unexpected code %d, expected %d", rr.Code, tc.wantCode)
			}

			if rr.Code == 200 {
				if rr.Body == nil {
					t.Fatal("userInfoHandler() returned no data")
				}

				var gotResults map[string]any
				if err = json.Unmarshal(rr.Body.Bytes(), &gotResults); err != nil {
					t.Fatalf("Failed to parse json data returned from userInfoHandler() %v", err)
				}

				if !reflect.DeepEqual(tc.wantResults, gotResults) {
					t.Fatalf("userInfoHandler() expected %v, got %v", tc.wantResults, gotResults)
				}
			}
		})
	}
}
