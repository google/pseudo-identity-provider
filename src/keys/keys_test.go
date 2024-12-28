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
)

func TestNoneSigner(t *testing.T) {
	noneSigner := NoneSigner{}

	sig, err := noneSigner.Sign([]byte{'a', 'b', 'c'}, nil)
	if err != nil {
		t.Errorf("unexpected error from nonesigner.Sign %v", err)
	}

	if sig != nil {
		t.Errorf("unexpected result from nonesigner.Sign %v", sig)
	}
}

func TestKeys(t *testing.T) {
	if err := SetupKeys(); err != nil {
		t.Fatalf("unexpected error from SetupKeys %v", err)
	}

	cases := []struct {
		alg     string
		algName string
		keyType string
	}{
		{
			alg:     "RSA",
			algName: "RS256",
			keyType: "RSA",
		},
		{
			alg:     "ES256",
			algName: "ES256",
			keyType: "EC",
		},
		{
			alg:     "ES384",
			algName: "ES384",
			keyType: "EC",
		},
		{
			alg:     "ES512",
			algName: "ES512",
			keyType: "EC",
		},
	}

	for _, tc := range cases {
		key := GetKey(tc.alg, false)
		if key == nil {
			t.Errorf("key not found %q", tc.alg)
			continue
		}

        if key.Jwk.Algorithm() != tc.algName {
            t.Errorf("unexpected algorithm from GetKey(%q) %q; want %q", tc.alg, key.Jwk.Algorithm(), tc.alg)
            }

		if key.Jwk.KeyType().String() != tc.keyType {
			t.Errorf("unexpected algorithm from GetKey(%q) %q; want %q", tc.alg, key.Jwk.KeyType().String(), tc.keyType)
		}
	}
}
