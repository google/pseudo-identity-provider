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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKeysHandler(t *testing.T) {
	config.SetGlobalConfig(&config.DefaultConfig)
	wantedKIDs := map[string]bool{
		"customidprsa":        true,
		"customidpecdsaP-256": true,
		"customidpecdsaP-384": true,
		"customidpecdsaP-521": true,
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(keyHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("keyHandler() returned unexpected failure code %d", rr.Code)
	}

	if rr.Body == nil {
		t.Fatal("keyHandler() returned no data")
	}

	var gotResults map[string]any
	if err = json.Unmarshal(rr.Body.Bytes(), &gotResults); err != nil {
		t.Fatalf("Failed to parse json data returned from keyHandler() %v", err)
	}

	keys, ok := gotResults["keys"]
	if !ok {
		t.Fatal("keyHandler() result does not have keys entry")
	}

	keysSlice, ok := keys.([]any)
	if !ok {
		t.Fatalf("keyHandler() result has invalid key list type %T", keys)
	}

	for _, key := range keysSlice {
		keyMap, ok := key.(map[string]any)
		if !ok {
			t.Fatalf("keyHandler() result has invalid key type %T", key)
		}

		kid, ok := keyMap["kid"]
		if !ok {
			t.Errorf("keyHandler() entry does not have Key Id %v", keyMap)
			continue
		}

		kidString, ok := kid.(string)
		if !ok {
			t.Errorf("keyHandler() result has invalid Key Id type %T", kid)
		}

		_, ok = wantedKIDs[kidString]
		if !ok {
			t.Errorf("keyHandler() returned unexpeced Key Id %q", kid)
		}
		delete(wantedKIDs, kidString)
	}

	if len(wantedKIDs) != 0 {
		t.Errorf("keyHandler() did not return expected Key Ids %v", wantedKIDs)
	}
}
