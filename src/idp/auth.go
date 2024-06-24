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
	"crypto/subtle"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

var (
	userNameVar = "LOG_USERNAME"
	pwdHashVar  = "LOG_PASSWORD"
	limiter = rate.NewLimiter(rate.Every(time.Minute/10), 10)
)

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	expectedUsername := os.Getenv(userNameVar)
	expectedPasswordHash := os.Getenv(pwdHashVar)

	// Consider auth turned off if both the expected username and password are not configured.
	if expectedUsername == "" && expectedPasswordHash == "" {
		return true
	}

	u, p, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restricted\"")
		http.Error(w, "", http.StatusUnauthorized)
		return false
	}

	res := limiter.Reserve()
	if !res.OK() || res.Delay() > (0 * time.Second) {
		http.Error(w, "Too many failures, please wait", http.StatusUnauthorized)
		return false
	}

	if subtle.ConstantTimeCompare([]byte(u), []byte(expectedUsername)) != 1 {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restricted\"")
		http.Error(w, "", http.StatusUnauthorized)
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(expectedPasswordHash), []byte(p))
	if err != nil {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"restricted\"")
		http.Error(w, "", http.StatusUnauthorized)
		return false
	}

	// Release limiter token on success. Cancel in the past so this reservation token
	// actually gets canceled.
	res.CancelAt(time.Now().Add(-time.Minute))
	return true
}
