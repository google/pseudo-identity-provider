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
	"bytes"
	"crypto/rand"
	"customidp/session"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
)

var randMethod func(b []byte) (int, error) = rand.Read

// Get evaluates the Parameter returing a slice of strings.
func (p Parameter) Get(input *session.RequestInput) ([]string, error) {
	switch p.Action {
	case "passthrough":
		return getInputParam(p.ID, input), nil
	case "set":
		return evaluateTemplates(p.Values, input)
	case "omit":
		return nil, nil
	case "random":
		b := make([]byte, 32)
		_, err := randMethod(b)
		if err != nil {
			return nil, err
		}
		// Encode our bytes as a base64 encoded string using URLEncoding
		encoded := base64.URLEncoding.EncodeToString(b)
		return []string{encoded}, nil
	case "custom":
		return getCustomParamValue(p.CustomKey, input)
	}
	return nil, nil
}

// GetJSON evaluates the Parameter as a JSON type.
func (p Parameter) GetJSON(input *session.RequestInput) (any, error) {
	vals, err := p.Get(input)
	if err != nil {
		return nil, err
	}
	return getJSONType(vals, p.JSONType)
}

// GetDefaultValue get a value for unconfigured parameters.
func (p Parameter) GetDefaultValue(input *session.RequestInput) []string {
	return getDefaultValue(p.ID, p.Action, input)
}

// getJSONType evaluates a value array as the requested type.
func getJSONType(vals []string, jsonType string) (any, error) {
	if len(vals) == 0 {
		return nil, nil
	}

	switch jsonType {
	case "string":
		return vals[0], nil
	case "array":
		return vals, nil
	case "number":
		num, err := strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a number %v", vals[0], err)
		}
		return num, nil
	case "boolean":
		truthy, err := strconv.ParseBool(vals[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as a boolean %v", vals[0], err)
		}
		return truthy, nil
	case "object":
		obj := make(map[string]any)
		if err := json.Unmarshal([]byte(vals[0]), &obj); err != nil {
			return nil, fmt.Errorf("failed to parse %q as a JSON Object %v", vals[0], err)
		}
		return obj, nil
	}

	// Default to string type.
	return vals[0], nil
}

// evaluateTemplates computes the golang template against the RequestInput.
func evaluateTemplates(templates []string, input *session.RequestInput) ([]string, error) {
	r := []string{}
	for _, temp := range templates {
		t, err := template.New("").Parse(temp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %v", err)
		}

		var b bytes.Buffer
		if err := t.Execute(&b, input); err != nil {
			return nil, fmt.Errorf("failed to execute template %v", err)
		}
		if b.String() != "" {
			r = append(r, b.String())
		}
	}
	return r, nil
}

// getInputParam gets the requested value from the input.
// If prioritized URL params over form params.
func getInputParam(id string, input *session.RequestInput) []string {
	val, ok := input.URLParams[id]
	if ok {
		return val
	}

	val, ok = input.FormParams[id]
	if ok {
		return val
	}

	return nil
}

// getDefaultValue get a value for unconfigured parameters.
func getDefaultValue(id string, action string, input *session.RequestInput) []string {
	switch action {
	case "passthrough":
		return getInputParam(id, input)
	case "omit":
		return nil
	}
	return nil
}
