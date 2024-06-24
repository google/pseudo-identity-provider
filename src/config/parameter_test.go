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
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"customidp/session"
)

func TestGet(t *testing.T) {
	RegisterCustomParam(
		"test_custom_param",
		func(input *session.RequestInput, config *Config) ([]string, error) {
			return []string{input.Domain}, nil
		})

	randMethod = func(b []byte) (int, error) {
		if len(b) > 6 {
			copy(b, []byte("random"))
			return 6, nil
		}

		return 0, fmt.Errorf("buffer too small")
	}

	cases := []struct {
		param Parameter
		input *session.RequestInput
		want  []string
	}{
		{
			param: Parameter{
				ID:     "test_passthrough",
				Action: "passthrough",
			},
			input: &session.RequestInput{
				URLParams: url.Values{
					"test_passthrough": []string{"test"},
				},
			},
			want: []string{"test"},
		},
		{
			param: Parameter{
				ID:     "test_passthrough_form",
				Action: "passthrough",
			},
			input: &session.RequestInput{
				FormParams: url.Values{
					"test_passthrough_form": []string{"test"},
				},
			},
			want: []string{"test"},
		},
		{
			param: Parameter{
				ID:     "test_set",
				Action: "set",
				Values: []string{"test1", "test2"},
			},
			want: []string{"test1", "test2"},
		},
		{
			param: Parameter{
				ID:     "test_set_template",
				Action: "set",
				Values: []string{"https://{{.Domain}}/oauth2/auth"},
			},
			input: &session.RequestInput{
				Domain: "test.com",
			},
			want: []string{"https://test.com/oauth2/auth"},
		},
		{
			param: Parameter{
				ID:     "test_omit",
				Action: "omit",
			},
			want: nil,
		},
		{
			param: Parameter{
				ID:     "test_random",
				Action: "random",
			},
			want: []string{"cmFuZG9tAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="},
		},
		{
			param: Parameter{
				ID:        "test_custom",
				Action:    "custom",
				CustomKey: "test_custom_param",
			},
			input: &session.RequestInput{
				Domain: "test.com",
			},
			want: []string{"test.com"},
		},
	}

	for _, tc := range cases {
		got, err := tc.param.Get(tc.input)
		if err != nil {
			t.Errorf("param.Get() for parameter %q failed %v", tc.param.ID, err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("param.Get() for parameter %q: expected %v, got %v", tc.param.ID, tc.want, got)
		}
	}
}

func TestGetJSON(t *testing.T) {
	cases := []struct {
		param Parameter
		input *session.RequestInput
		want  any
		err   error
	}{
		{
			param: Parameter{
				ID:       "test_string",
				Action:   "set",
				Values:   []string{"test1"},
				JSONType: "string",
			},
			want: "test1",
			err:  nil,
		},
		{
			param: Parameter{
				ID:       "test_array",
				Action:   "set",
				Values:   []string{"test1", "test2"},
				JSONType: "array",
			},
			want: []string{"test1", "test2"},
			err:  nil,
		},
		{
			param: Parameter{
				ID:       "test_number",
				Action:   "set",
				Values:   []string{"10"},
				JSONType: "number",
			},
			want: 10,
			err:  nil,
		},
		{
			param: Parameter{
				ID:       "test_number_invalid",
				Action:   "set",
				Values:   []string{"10abc"},
				JSONType: "number",
			},
			err: fmt.Errorf(""),
		},
		{
			param: Parameter{
				ID:       "test_boolean",
				Action:   "set",
				Values:   []string{"true"},
				JSONType: "boolean",
			},
			want: true,
			err:  nil,
		},
		{
			param: Parameter{
				ID:       "test_boolean_invalid",
				Action:   "set",
				Values:   []string{"nottrue"},
				JSONType: "boolean",
			},
			err: fmt.Errorf(""),
		},
		{
			param: Parameter{
				ID:       "test_object",
				Action:   "set",
				Values:   []string{`{"testobj": { "field": "val"}}`},
				JSONType: "object",
			},
			want: map[string]any{
				"testobj": map[string]any{
					"field": "val",
				},
			},
			err: nil,
		},
		{
			param: Parameter{
				ID:       "test_object_invalid",
				Action:   "set",
				Values:   []string{`"testobj": { "field": "val"}`},
				JSONType: "object",
			},
			err: fmt.Errorf(""),
		},
		{
			param: Parameter{
				ID:       "test_empty",
				Action:   "set",
				JSONType: "string",
			},
			want: nil,
			err:  nil,
		},
		{
			param: Parameter{
				ID:     "test_default",
				Action: "set",
				Values: []string{"test1"},
			},
			want: "test1",
			err:  nil,
		},
	}

	for _, tc := range cases {
		got, err := tc.param.GetJSON(tc.input)
		if err != nil && tc.err == nil {
			t.Errorf("param.GetJSON() failed with unexpected error: %v", err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("param.GetJSON() for parameter %q: expected %v, got %v", tc.param.ID, tc.want, got)
		}
	}
}

func TestGetDefaultValue(t *testing.T) {
	cases := []struct {
		param Parameter
		input *session.RequestInput
		want  []string
	}{
		{
			param: Parameter{
				ID:     "test_passthrough",
				Action: "passthrough",
			},
			input: &session.RequestInput{
				URLParams: url.Values{
					"test_passthrough": []string{"test"},
				},
			},
			want: []string{"test"},
		},
		{
			param: Parameter{
				ID:     "test_passthrough_form",
				Action: "passthrough",
			},
			input: &session.RequestInput{
				FormParams: url.Values{
					"test_passthrough_form": []string{"test"},
				},
			},
			want: []string{"test"},
		},
		{
			param: Parameter{
				ID:     "test_omit",
				Action: "omit",
			},
			want: nil,
		},
	}

	for _, tc := range cases {
		got := tc.param.GetDefaultValue(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("param.GetDefaultValue() for parameter %q: expected %v, got %v", tc.param.ID, tc.want, got)
		}
	}
}
