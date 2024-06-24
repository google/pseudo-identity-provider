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
	sessionmgmt "customidp/session"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine/v2"
)

// errorResponse returns an error based on configuration.
func errorResponse(w http.ResponseWriter, r *http.Request, e *config.Error) {
	http.Error(w, e.ErrorContent, e.ErrorCode)
}

// blockResponse blocks the return of a handler for a long time.
func blockResponse(w http.ResponseWriter) {
	time.Sleep(10 * time.Minute)
	http.Error(w, "", http.StatusGatewayTimeout)
}

// jsonResponse builds a JSON formated response from configured Parameter values.
func jsonResponse(w http.ResponseWriter, input *sessionmgmt.RequestInput, parameters []config.Parameter) {
	content := map[string]any{}
	for _, configParam := range parameters {
		jsonVal, err := configParam.GetJSON(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if jsonVal != nil {
			content[configParam.ID] = jsonVal
		}
	}

	resp, err := json.Marshal(content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal content %v", err), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(resp))
}

// getDomain constructs the instances domain name from AppEngine config.
func getDomain(r *http.Request) string {
	ctx := appengine.NewContext(r)
	if !appengine.IsAppEngine() {
		return r.Host
	}

	if appengine.ModuleName(ctx) == "default" {
		return appengine.DefaultVersionHostname(ctx)
	}

	return appengine.ModuleName(ctx) + "-dot-" + appengine.DefaultVersionHostname(ctx)
}

// getInputData extracts important info from the request.
func getInputData(r *http.Request) *sessionmgmt.RequestInput {
	var session sessionmgmt.Session
	if err := r.ParseForm(); err != nil {
		logError(fmt.Sprintf("failed making input struct: %v", err), r)
	}

	// If code is an input parameter (token endpoint). Load existing session state.
	if r.Form.Get("code") != "" {
		var err error
		session, err = sessionmgmt.GetSession(r.Form.Get("code"))
		if err != nil {
			logError(fmt.Sprintf("unexpected code: %v", err), r)
		}
	}

	return &sessionmgmt.RequestInput{
		HTTPMethod: r.Method,
		Path:       r.URL.Path,
		Proto:      r.Header.Get("X-Forwarded-Proto"),
		Headers:    r.Header,
		URLParams:  r.URL.Query(),
		FormParams: r.PostForm,
		Domain:     getDomain(r),
		Session:    &session,
		Time:       time.Now()}
}
