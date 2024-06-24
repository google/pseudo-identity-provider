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
	"customidp/session"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestTokenHandler(t *testing.T) {
	config.RegisterCustomParam(
		"test_custom_param",
		func(input *session.RequestInput, config *config.Config) ([]string, error) {
			return []string{"https://customidp.com/callback"}, nil
		})

	cases := []struct {
		title       string
		config      *config.Config
		wantCode    int
		wantResults map[string]any
		clearValues bool
	}{
		{
			title:    "Default config",
			wantCode: 200,
			wantResults: map[string]any{
				"access_token":  "",
				"id_token":      "",
				"refresh_token": "",
				"token_type":    "",
				"expires_in":    "",
			},
			clearValues: true,
		},
		{
			title: "Error Response",
			config: &config.Config{
				TokenAction: config.TokenAction{
					Action: "error",
					Error: config.Error{
						ErrorCode:    404,
						ErrorContent: "Not Found",
					},
				},
			},
			wantCode: 404,
		},
		{
			title: "Use session state",
			config: &config.Config{
				TokenAction: config.TokenAction{
					Action: "respond",
					Respond: config.TokenRespond{
						Parameters: []config.Parameter{
							{ID: "id_token", Action: "set", Values: []string{"{{.Session.Code}}"}, JSONType: "string"},
						},
					},
				},
			},
			wantCode: 200,
			wantResults: map[string]any{
				"id_token": "randomval",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			if tc.config != nil {
				config.SetGlobalConfig(tc.config)
			} else {
				config.SetGlobalConfig(&config.DefaultConfig)
			}

			session.CreateSession(&session.RequestInput{}, url.Values{"code": []string{"randomval"}})

			form := `code=randomval&
			client_id=your_client_id&
			client_secret=your_client_secret&
			redirect_uri=https%3A//idpclient.idp/code&
			grant_type=authorization_code`

			req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(form)))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(tokenHandler)

			handler.ServeHTTP(rr, req)

			if tc.wantCode != rr.Code {
				t.Errorf("tokenHandler() returned %d rather than expected %d", rr.Code, tc.wantCode)
			}

			if rr.Code == 200 {
				if rr.Body == nil {
					t.Fatal("tokenHandler() returned no data")
				}

				var gotResults map[string]any
				if err = json.Unmarshal(rr.Body.Bytes(), &gotResults); err != nil {
					t.Fatalf("Failed to parse json data returned from discHandler() %v", err)
				}

				if tc.clearValues {
					for k := range gotResults {
						gotResults[k] = ""
					}
				}

				if !reflect.DeepEqual(tc.wantResults, gotResults) {
					t.Fatalf("discHandler() expected %v, got %v", tc.wantResults, gotResults)
				}
			}
		})
	}
}
