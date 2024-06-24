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

// Package idp implements the main request Handler logic for the IdP.
package idp

import (
	"customidp/keys"
	"html/template"
	"net/http"
)

// InitHandlers configures all HTTP request handlers for the IdP.
func InitHandlers(staticTemplates *template.Template, staticDir string) error {
	if err := keys.SetupKeys(); err != nil {
		return err
	}

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", staticWithCSP(fs, staticTemplates))

	http.HandleFunc("/log", logHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/configschema", configSchemaHandler)
	http.HandleFunc("/.well-known/openid-configuration", respLogHandler(discHandler))
	http.HandleFunc("/.well-known/jwks.json", respLogHandler(keyHandler))
	http.HandleFunc("/oauth2/auth", respLogHandler(authHandler))
	http.HandleFunc("/oauth2/token", respLogHandler(tokenHandler))
	http.HandleFunc("/oauth2/userinfo", respLogHandler(userInfoHandler))
	return nil
}
