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

// userInfoHandler takes a token and returns associated user information.
func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	action := config.GetGlobalConfig().UserInfoAction
	input := getInputData(r)
	addRequestLogEntry(input, action.Action)

	switch action.Action {
	case "respond":
		userInfoRespond(w, input)
	case "error":
		errorResponse(w, r, &action.Error)
	case "block":
		blockResponse(w)
	}
}

// userInfoRespond responds with JSON content as configured.
func userInfoRespond(w http.ResponseWriter, input *sessionmgmt.RequestInput) {
	c := config.GetGlobalConfig().UserInfoAction.Respond
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	jsonResponse(w, input, c.Parameters)
}
