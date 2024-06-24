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
	"fmt"
)

var customMethodMap map[string]func(input *sessionmgmt.RequestInput, config *Config) ([]string, error)

// RegisterCustomParam registers a custom parameter evaluation function.
func RegisterCustomParam(id string, f func(input *sessionmgmt.RequestInput, config *Config) ([]string, error)) {
	if customMethodMap == nil {
		customMethodMap = make(map[string]func(input *sessionmgmt.RequestInput, config *Config) ([]string, error))
	}

	customMethodMap[id] = f
}

// getCustomParamValue calls the configured custom parameter evaluation function.
func getCustomParamValue(key string, input *sessionmgmt.RequestInput) ([]string, error) {
	if customMethodMap == nil {
		return nil, fmt.Errorf("custom method key %q is not defined", key)
	}

	f, ok := customMethodMap[key]
	if !ok {
		return nil, fmt.Errorf("custom method key %q is not defined", key)
	}

	return f(input, GetGlobalConfig())
}
