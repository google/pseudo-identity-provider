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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// staticWithCSP serves static files and evaluates HTML files as GoLang templates to inject
// CSP nonces.
func staticWithCSP(fs http.Handler, staticTemplates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nonce, err := getNonce()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Security-Policy",
			fmt.Sprintf(`default-src 'self'; style-src 'self' 'nonce-%s';
			 	script-src 'self' 'nonce-%s'; font-src https://fonts.gstatic.com/;`, nonce, nonce))

		// Serve static files, if it isn't an expected file path, assume it is a Angular route and
		// load index.html.
		file := path.Base(r.URL.Path)
		if path.Ext(file) == ".js" || path.Ext(file) == ".css" || path.Ext(file) == ".ico" {
			fs.ServeHTTP(w, r)
			return
		}

		template := staticTemplates.Lookup("index.html")
		if template != nil {
			nonce := struct {
				CspNonce string
			}{
				CspNonce: nonce,
			}

			template.Execute(w, nonce)
			return
		}
	}
}

// getNonce returns a random nonce for CSP.
func getNonce() (string, error) {
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(nonce), nil
}
