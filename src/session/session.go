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

// Package session tracks auth session across calls.
package session

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Session State tracking.
type Session struct {
	// OAuth code value from Authorization Endpoint.
	Code                string

	// Nonce value from client if provided.
	Nonce               string

	// PKCE code challenge value if provided.
	CodeChallenge       string

	// PKCE code challenge method if provided.
	CodeChallengeMethod string

	// OAuth Client ID.
	ClientID            string

	// The client's redirect URI specified at the Authorization Endpoint.
	RedirectURI         string
}

// RequestInput tracks request state and can be use in Parameter evaluation templates.
type RequestInput struct {
	// Domain name of the IdP Server.
	Domain     string

	// The HTTP Method Used by the caller. GET, POST, etc.
	HTTPMethod string

	// The URL Path used in the call.
	Path       string

	// HTTPS if the call is over TLS, HTTP otherwise.
	Proto      string

	// HTTP Header values.
	Headers    http.Header

	// URL Query parameter values. Can be indexed like a Map.
	URLParams  url.Values

	// POST form parameters if any. Can be indexed like a Map.
	FormParams url.Values

	// The session structure for Token endpoint calls if any is found.
	Session    *Session

	// Call timestamp.
	Time       time.Time
}

// Global map for tracking sessions.
var sessions map[string]Session
var sessionsMutex sync.Mutex

// CreateSession adds a new session and pulls session state from the Request input.
func CreateSession(input *RequestInput, updatedParams url.Values) {
	// Key the session by the authorization code returned by the IdP.
	code := updatedParams.Get("code")
	if code == "" {
		// If there is no code, we have nothing to later key the session by in token endpoint.
		// This likely indicates Implicit mode where tracking the session across calls isn't necessary.
		return
	}

	// Get the updated redirect uri and if none, the input redirect uri.
	redirectURI := updatedParams.Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = input.URLParams.Get("redirect_uri")
	}

	session := Session{
		ClientID:            input.URLParams.Get("client_id"),
		Nonce:               input.URLParams.Get("nonce"),
		CodeChallenge:       input.URLParams.Get("code_challenge"),
		CodeChallengeMethod: input.URLParams.Get("code_challenge_method"),
		RedirectURI:         redirectURI,
		Code:                code,
	}

	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	if sessions == nil {
		sessions = make(map[string]Session)
	}
	sessions[code] = session
}

// GetSession returns the Session by code key.
func GetSession(code string) (Session, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	session, ok := sessions[code]
	if !ok {
		return Session{}, fmt.Errorf("no session found")
	}

	return session, nil
}
