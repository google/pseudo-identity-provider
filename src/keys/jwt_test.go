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

package keys

import (
	"testing"

	"github.com/lestrrat-go/jwx/jwt"
)

func TestSignToken(t *testing.T) {
	if err := SetupKeys(); err != nil {
		t.Fatalf("unexpected error from SetupKeys %v", err)
	}

	token := jwt.New()
	token.Set("sub", "testsub")

	cases := []struct {
		sigAlg string
	}{
		{
			sigAlg: "RS256",
		},
		{
			sigAlg: "RS384",
		},
		{
			sigAlg: "RS512",
		},
		{
			sigAlg: "ES256",
		},
		{
			sigAlg: "ES384",
		},
		{
			sigAlg: "ES512",
		},
		{
			sigAlg: "HS256",
		},
		{
			sigAlg: "none",
		},
	}

	for _, tc := range cases {
		signedToken, err := SignToken(tc.sigAlg, token, false)
		if err != nil {
			t.Errorf("SignToken(%q) failed: %v", tc.sigAlg, err)
		}

		_, err = jwt.Parse([]byte(signedToken))
		if err != nil {
			t.Errorf("SignToken(%q) returned unparsable token: %v", tc.sigAlg, err)
		}
	}
}
