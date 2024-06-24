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
	"net/http"
)

// discHandler returns OIDC Discovery doc.
func discHandler(w http.ResponseWriter, r *http.Request) {
	action := config.GetGlobalConfig().DiscoveryAction
	input := getInputData(r)
	addRequestLogEntry(input, action.Action)

	switch action.Action {
	case "respond":
		discRespond(w, r, input)
	case "error":
		errorResponse(w, r, &action.Error)
	case "block":
		blockResponse(w)
	}
}

// discRespond responds with the generated discovery doc.
func discRespond(w http.ResponseWriter, r *http.Request, input *sessionmgmt.RequestInput) {
	respond := config.GetGlobalConfig().DiscoveryAction.Respond

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	jsonResponse(w, input, respond.Parameters)
}
