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

package config

import (
	"customidp/keys"
	"customidp/session"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func TestGenerateToken(t *testing.T) {
	if err := keys.SetupKeys(); err != nil {
		t.Fatal(err)
	}

	input := &session.RequestInput{
		Domain: "text.com",
	}

	expectedIss := "https://test.com"
	expectedSub := "12345abcde"

	config := &Config{
		IDTokenConfig: IDTokenConfig{
			Algorithm: "RS256",
			Claims: []Claim{
				{ID: "iss", Values: []string{"https://{{.Domain}}"}, JSONType: "string"},
				{ID: "sub", Values: []string{expectedSub}, JSONType: "string"},
			},
			RemoveSignature: false,
			UseWrongKey:     false,
		},
	}

	got, err := GenerateToken(input, config)
	if err != nil {
		t.Errorf("GenerateToken() failed: %v", err)
	}

	pubKey, err := keys.GetKey("RSA", false).Jwk.PublicKey()
	if err != nil {
		t.Errorf("Failedto get RSA key: %v", err)
	}

	_, err = jwt.ParseString(
		got[0],
		jwt.WithVerify(jwa.RS256, pubKey),
		jwt.WithClaimValue("iss", expectedIss),
		jwt.WithClaimValue("sub", expectedSub))
	if err != nil {
		t.Errorf("GenerateToken() created unparsable token: %v", err)
	}
}
