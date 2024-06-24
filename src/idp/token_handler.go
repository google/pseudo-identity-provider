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

// tokenHandler takes action for the Token Endpoint based on config.
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// OAuth Spec says clients must use POST, however we won't enforce that
	// here. We will just log the method along with other request info.
	action := config.GetGlobalConfig().TokenAction
	input := getInputData(r)
	addRequestLogEntry(input, action.Action)

	switch action.Action {
	case "respond":
		tokenRespond(w, input)
	case "error":
		errorResponse(w, r, &action.Error)
	case "block":
		blockResponse(w)
	}
}

// tokenRespond responds with jsonContent based on configuration.
// The "signed_token_id" custom method can be used to create valid
// signed tokens.
func tokenRespond(w http.ResponseWriter, input *sessionmgmt.RequestInput) {
	c := config.GetGlobalConfig().TokenAction.Respond
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	jsonResponse(w, input, c.Parameters)
}
