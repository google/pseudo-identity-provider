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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

var (
	testDefaultUsername = "test"
	testDefaultPassword = "testpwd"
)

func setupCreds(t *testing.T) {
	expectedPwdhash, err := bcrypt.GenerateFromPassword([]byte(testDefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv(userNameVar, testDefaultUsername)
	os.Setenv(pwdHashVar, string(expectedPwdhash))
}

func TestCheckAuth(t *testing.T) {
	setupCreds(t)

	cases := []struct {
		title    string
		username string
		password string
		want     bool
	}{
		{
			title:    "Valid username/password",
			username: testDefaultUsername,
			password: testDefaultPassword,
			want:     true,
		},
		{
			title:    "Invalid username",
			username: "wrongusername",
			password: testDefaultPassword,
			want:     false,
		},
		{
			title:    "Invalid password",
			username: testDefaultPassword,
			password: "wrongpassword",
			want:     false,
		},
	}

	for _, tc := range cases {
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.SetBasicAuth(tc.username, tc.password)

		w := httptest.NewRecorder()

		got := checkAuth(w, r)
		if got != tc.want {
			t.Errorf("expected %t, got %t", tc.want, got)
		}
	}
}
