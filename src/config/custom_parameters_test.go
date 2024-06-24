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

package config

import (
	sessionmgmt "customidp/session"
	"reflect"
	"testing"
)

func TestGetCustomParamValue(t *testing.T) {
	// Verify that the custom parameter method is passed the expected input and config.
	RegisterCustomParam(
		"test_custom_param",
		func(input *sessionmgmt.RequestInput, config *Config) ([]string, error) {
			return []string{input.Domain, config.AuthAction.Action}, nil
		})

	input := &sessionmgmt.RequestInput{
		Domain: "test.com",
	}

	c := &Config{
		AuthAction: AuthAction{
			Action: "reply",
		},
	}
	SetGlobalConfig(c)

	want := []string{"test.com", "reply"}
	got, _ := getCustomParamValue("test_custom_param", input)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
