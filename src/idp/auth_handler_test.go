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
	"customidp/session"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	config.RegisterCustomParam(
		"test_custom_param",
		func(input *session.RequestInput, config *config.Config) ([]string, error) {
			return []string{"https://customidp.com/callback"}, nil
		})

	cases := []struct {
		title            string
		config           *config.Config
		url              string
		urlParams        url.Values
		wantCode         int
		wantRedirectURL  string
		wantURLParams    url.Values
		wantHashFragment bool
	}{
		{
			title: "Default config redirect",
			url:   "/oauth2/auth",
			urlParams: url.Values{
				"redirect_uri": {"https://localhost:8080/callback"},
				"state":        {"randstate"},
			},
			wantCode:        302,
			wantRedirectURL: "https://localhost:8080/callback",
			wantURLParams: url.Values{
				"code":  nil, // code is random so ignore the value.
				"state": {"randstate"},
			},
		},
		{
			title: "Implicit Flow with Hash",
			config: &config.Config{
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
			},
			url: "/oauth2/auth",
			urlParams: url.Values{
				"redirect_uri": {"https://localhost:8080/callback"},
				"state":        {"randstate"},
			},
			wantCode:        302,
			wantRedirectURL: "https://localhost:8080/callback",
			wantURLParams: url.Values{
				"access_token": nil, // code is random so ignore the value.
				"token_type":   {"Bearer"},
				"state":        {"randstate"},
			},
			wantHashFragment: true,
		},
		{
			title: "Custom redirect",
			config: &config.Config{
				AuthAction: config.AuthAction{
					Action: "redirect",
					Redirect: config.AuthRedirect{
						DefaultParamAction: "passthrough",
						Parameters: []config.Parameter{
							{ID: "client_id", Action: "set", Values: []string{"testid"}},
						},
						RedirectTarget: config.RedirectTarget{
							UseCustomRedirectURI: true,
							Target:               "https://anotheridp.com/callback",
						},
					},
				},
			},
			url: "/oauth2/auth",
			urlParams: url.Values{
				"redirect_uri": {"https://localhost:8080/callback"},
				"state":        {"randstate"},
			},
			wantCode:        302,
			wantRedirectURL: "https://anotheridp.com/callback",
			wantURLParams: url.Values{
				"client_id":    {"testid"},
				"state":        {"randstate"},
				"redirect_uri": {"https://localhost:8080/callback"},
			},
		},
		{
			title: "Custom redirect custom parameter",
			config: &config.Config{
				AuthAction: config.AuthAction{
					Action: "redirect",
					Redirect: config.AuthRedirect{
						DefaultParamAction: "passthrough",
						Parameters: []config.Parameter{
							{ID: "client_id", Action: "set", Values: []string{"testid"}},
						},
						RedirectTarget: config.RedirectTarget{
							UseCustomRedirectURI: true,
							CustomKey:            "test_custom_param",
						},
					},
				},
			},
			url: "/oauth2/auth",
			urlParams: url.Values{
				"redirect_uri": {"https://localhost:8080/callback"},
				"state":        {"randstate"},
			},
			wantCode:        302,
			wantRedirectURL: "https://customidp.com/callback",
			wantURLParams: url.Values{
				"client_id":    {"testid"},
				"state":        {"randstate"},
				"redirect_uri": {"https://localhost:8080/callback"},
			},
		},
		{
			title: "Omit parameters",
			config: &config.Config{
				AuthAction: config.AuthAction{
					Action: "redirect",
					Redirect: config.AuthRedirect{
						DefaultParamAction: "omit",
						Parameters: []config.Parameter{
							{ID: "client_id", Action: "set", Values: []string{"testid"}},
						},
					},
				},
			},
			url: "/oauth2/auth",
			urlParams: url.Values{
				"redirect_uri": {"https://localhost:8080/callback"},
				"state":        {"randstate"},
			},
			wantCode:        302,
			wantRedirectURL: "https://localhost:8080/callback",
			wantURLParams: url.Values{
				"client_id": {"testid"},
			},
		},
		{
			title: "Error Response",
			config: &config.Config{
				AuthAction: config.AuthAction{
					Action: "error",
					Error: config.Error{
						ErrorCode:    404,
						ErrorContent: "Not Found",
					},
				},
			},
			url:       "/oauth2/auth",
			urlParams: url.Values{"redirect_uri": {"https://localhost:8080/callback"}},
			wantCode:  404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			if tc.config != nil {
				config.SetGlobalConfig(tc.config)
			} else {
				config.SetGlobalConfig(&config.DefaultConfig)
			}

			req, err := http.NewRequest("GET", tc.url+"?"+tc.urlParams.Encode(), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(authHandler)

			handler.ServeHTTP(rr, req)

			if tc.wantCode != rr.Code {
				t.Errorf("authHandler() returned %d rather than expected %d", rr.Code, tc.wantCode)
			}

			if tc.wantRedirectURL != "" {
				locations := rr.Result().Header["Location"]
				if len(locations) != 1 {
					t.Errorf("unexpected number of redirect uris returned %v", locations)
				}

				gotURL, err := url.Parse(locations[0])
				if err != nil {
					t.Errorf("expected redirect URL but it failed to parse %v", err)
				}

				if !strings.HasPrefix(gotURL.String(), tc.wantRedirectURL) {
					t.Errorf("expected redirect URL starting with %q but got %q", tc.wantRedirectURL, gotURL.String())
				}

				params := gotURL.Query()
				if tc.wantHashFragment {
					params, err = url.ParseQuery(gotURL.Fragment)
					if err != nil {
						t.Errorf("failed to parse fragment %v", err)
					}
				}

				// Check that we got the parameters we expected.
				for wantID, wantVal := range tc.wantURLParams {
					gotVal, ok := params[wantID]
					if !ok {
						t.Errorf("expected URL parameter %q but it was not returned", wantID)
					}

					if wantVal != nil {
						if !reflect.DeepEqual(wantVal, gotVal) {
							t.Errorf("authHandler() for parameter %q: expected %v, got %v", wantID, wantVal, gotVal)
						}
					}
				}

				// And didn't get any that we did not expect.
				for gotID := range params {
					_, ok := tc.wantURLParams[gotID]
					if !ok {
						t.Errorf("unexpected URL parameter %q", gotID)
					}
				}
			}
		})
	}
}
