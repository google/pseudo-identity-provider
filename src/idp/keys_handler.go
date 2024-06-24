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
	"customidp/keys"
	"fmt"
	"net/http"
)

// keyHandler returns siging keys as a JWK Set in JSON.
func keyHandler(w http.ResponseWriter, r *http.Request) {
	input := getInputData(r)
	addRequestLogEntry(input, "")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json, err := keys.GetJSONKeySet()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get key set %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, json)
}
