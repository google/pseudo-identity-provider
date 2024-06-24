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

package session

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestSessionStorage(t *testing.T) {
	cases := []struct {
		title				string
		code        string
		redirectURI string
		clientID    string
		session     Session
		err         error
	}{
		{
			title:       "Session Found",
			code:        "random",
			redirectURI: "https://test.com",
			clientID:    "testid",
			session: Session{
				Code:        "random",
				RedirectURI: "https://test.com",
				ClientID:    "testid",
			},
			err: nil,
		},
		{
			title:       "Session Not Found",
			code:        "",
			redirectURI: "https://test.com",
			clientID:    "testid",
			err:         fmt.Errorf("no session found"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			input := &RequestInput{
				URLParams: url.Values{
					"client_id": []string{tc.clientID},
				},
			}

			updatedParams := url.Values{
				"redirect_uri": []string{tc.redirectURI},
				"code":         []string{tc.code},
			}

			CreateSession(input, updatedParams)
			session, err := GetSession(tc.code)
			if err != nil && tc.err == nil {
				t.Errorf("GetSession() failed with unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tc.session, session) {
				t.Errorf("expected %v, got %v", tc.session, session)
			}
		})
	}
}
